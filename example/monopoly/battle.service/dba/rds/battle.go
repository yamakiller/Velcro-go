package rds

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lithammer/shortuuid/v4"
	"github.com/yamakiller/velcro-go/example/monopoly/battle.service/errs"
	"github.com/yamakiller/velcro-go/example/monopoly/pub/rdsconst"
	"github.com/yamakiller/velcro-go/network"
)

//1.创建对战区 对战区长时不开始且没有任何人加入，将在10分钟后结束
//2.删除对战区
//3.加入对战区
//4.离开对战区
//5.查询对战区列表
//6.自动加入某对战区
//7.对战始流传: 7-1.通知主机提交UDP(既:Nat信息)
//				 7-2.收到Nat信息回复客户端，已收到.(不需要再发送)
//				 7-3.将Nat信息发送给参与者
//				 7-4.通知各方开始游戏:
//									7-4-1.如果通知失败,通知主机对战取消,对战区取消开始状态
//									7-4-2.修改对战区状态,修改对战区倒计时时间

func CreateBattleSpace(ctx context.Context,
	uid string,
	clientId *network.ClientID,
	mapURi string,
	max_count int32) (string, error) {

	mutex := sync.NewMutex(rdsconst.GetPlayerLockKey(uid))
	if err := mutex.Lock(); err != nil {
		return "", err
	}

	results, err := client.HMGet(ctx, rdsconst.GetPlayerOnlineDataKey(uid),
		rdsconst.PlayerMapClientIcon,
		rdsconst.PlayerMapClientBattleSpaceId,
		rdsconst.PlayerMapClientDisplayName).Result()
	if err != nil {
		mutex.Unlock()
		return "", err
	}

	if len(results) != 3 {
		mutex.Unlock()
		return "", errs.ErrorPlayerOnlineDataLost
	}

	if results[1] != "" {
		mutex.Unlock()
		return "", errs.ErrorPlayerAlreadyInBattleSpace
	}

	pipe := client.TxPipeline()
	defer pipe.Close()

	pipe.Do(ctx, "MULTI")

	results[1] = shortuuid.New()

	battleSpace := map[string]string{
		rdsconst.BattleSpaceId:                    results[1].(string),
		rdsconst.PalyerMapClientIdAddress:         clientId.Address,
		rdsconst.PlayerMapClientIdId:              clientId.Id,
		rdsconst.BattleSpaceMapURi:                mapURi,
		rdsconst.BattleSpaceNatAddr:               "",
		rdsconst.BattleSpaceMasterUid:             uid,
		rdsconst.BattleSpaceMasterIcon:            results[0].(string),
		rdsconst.BattleSpaceMasterDisplay:         rdsconst.BattleSpaceStateReady,
		rdsconst.BattleSpacePlayerCount:           strconv.FormatInt(int64(max_count), 10),
		rdsconst.BattleSpacePlayerPos:             rdsconst.UpdateData(make([]string, max_count), 0, uid),
		rdsconst.GetBattleSpacePlayerDataKey(uid): fmt.Sprintf("%s&%s&%s&%s&%s&%s", uid, results[0].(string), results[2].(string), rdsconst.BattleSpaceStateReady, "0"),
		rdsconst.BattleSpaceTime:                  strconv.FormatInt(time.Now().UnixMilli(), 10),
		rdsconst.BattleSpaceState:                 rdsconst.BattleSpaceStateNomal,
	}

	pipe.HSet(ctx, rdsconst.GetPlayerOnlineDataKey(uid), rdsconst.PlayerMapClientBattleSpaceId, results[1])
	pipe.HMSet(ctx, rdsconst.GetBattleSpaceOnlineDataKey(results[1].(string)), battleSpace)
	pipe.RPush(ctx, rdsconst.BattleSpaceOnlinetable, results[1].(string))

	pipe.Do(ctx, "exec")

	_, err = pipe.Exec(ctx)
	if err != nil {
		pipe.Discard()
		mutex.Unlock()
		return "", err
	}

	mutex.Unlock()

	return results[0].(string), nil
}

func DeleteBattleSpace(ctx context.Context, clientId *network.ClientID) error {
	uid, err := client.Get(ctx, rdsconst.GetPlayerClientIDKey(clientId.ToString())).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		return err
	}

	player_mutex := sync.NewMutex(rdsconst.GetPlayerLockKey(uid))
	if err := player_mutex.Lock(); err != nil {
		return err
	}
	results, err := client.HMGet(ctx, rdsconst.GetPlayerOnlineDataKey(uid),
		rdsconst.PlayerMapClientBattleSpaceId).Result()
	if err != nil {
		player_mutex.Unlock()
		return err
	}
	if len(results) != 1 {
		player_mutex.Unlock()
		return errs.ErrorPlayerOnlineDataLost
	}
	if results[0] == "" {
		player_mutex.Unlock()
		return errs.ErrorPlayerIsNotInBattleSpace
	}

	spaceid := results[0].(string)

	space_mutex := sync.NewMutex(rdsconst.GetBattleSpaceOnlineDataKey(spaceid))
	if err := space_mutex.Lock(); err != nil {
		return err
	}

	spaceResults, err := client.HGetAll(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid)).Result()
	if err != nil {
		space_mutex.Unlock()
		player_mutex.Unlock()
		return err
	}

	if !clientId.Equal(&network.ClientID{Address: spaceResults[rdsconst.PalyerMapClientIdAddress],
		Id: spaceResults[rdsconst.PlayerMapClientIdId]}) {
		space_mutex.Unlock()
		player_mutex.Unlock()
		return errors.New("permissions lost")
	}

	pipe := client.TxPipeline()
	defer pipe.Close()

	pipe.Do(ctx, "MULTI")

	pipe.HSet(ctx, rdsconst.GetPlayerOnlineDataKey(uid), rdsconst.PlayerMapClientBattleSpaceId, "")
	pipe.HDel(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid))
	pipe.LRem(ctx, rdsconst.BattleSpaceOnlinetable, 1, spaceid)
	pipe.Do(ctx, "exec")

	_, err = pipe.Exec(ctx)
	if err != nil {
		pipe.Discard()
		space_mutex.Unlock()
		player_mutex.Unlock()
		return err
	}

	space_mutex.Unlock()
	player_mutex.Unlock()
	return nil
}

func EnterBattleSpace(ctx context.Context, spaceid string, clientId *network.ClientID) error {
	uid, err := client.Get(ctx, rdsconst.GetPlayerClientIDKey(clientId.ToString())).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		return err
	}

	player_mutex := sync.NewMutex(rdsconst.GetPlayerLockKey(uid))
	if err := player_mutex.Lock(); err != nil {
		return err
	}
	playerResults, err := client.HMGet(ctx, rdsconst.GetPlayerOnlineDataKey(uid),
		rdsconst.PlayerMapClientIcon,
		rdsconst.PlayerMapClientBattleSpaceId,
		rdsconst.PlayerMapClientDisplayName).Result()
	if err != nil {
		player_mutex.Unlock()
		return err
	}
	if len(playerResults) != 3 {
		player_mutex.Unlock()
		return errs.ErrorPlayerOnlineDataLost
	}
	if playerResults[1] != "" {
		player_mutex.Unlock()
		return errs.ErrorPlayerAlreadyInBattleSpace
	}

	space_mutex := sync.NewMutex(rdsconst.GetBattleSpaceOnlineDataKey(spaceid))
	if err := space_mutex.Lock(); err != nil {
		player_mutex.Unlock()
		return err
	}

	player_results, err := client.HMGet(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid), rdsconst.GetBattleSpacePlayerDataKey("")).Result()
	if err != nil {
		space_mutex.Unlock()
		player_mutex.Unlock()
		return err
	}
	for _, player := range player_results {
		if player == rdsconst.GetBattleSpacePlayerDataKey(uid) {
			return errs.ErrorPlayerAlreadyInBattleSpace
		}
	}

	max_count, err := client.HMGet(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid), rdsconst.GetBattleSpacePlayerDataKey("")).Result()
	if err != nil {
		space_mutex.Unlock()
		player_mutex.Unlock()
		return err
	}
	count, _ := strconv.Atoi(max_count[0].(string))
	if count <= len(player_results) {
		return errs.ErrorSpacePlayerIsFull
	}

	pos_result, err := client.HMGet(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid), rdsconst.GetBattleSpacePlayerDataKey("")).Result()
	if err != nil {
		space_mutex.Unlock()
		player_mutex.Unlock()
		return err
	}
	space_pos, index := enterBattleSpacePlayerPos(pos_result[0].(string), uid)
	pipe := client.TxPipeline()
	defer pipe.Close()

	pipe.Do(ctx, "MULTI")
	pipe.HSet(ctx, rdsconst.GetPlayerOnlineDataKey(uid), rdsconst.PlayerMapClientBattleSpaceId, spaceid)
	pipe.HSet(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid), rdsconst.BattleSpacePlayerPos, space_pos)
	pipe.HSet(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid),
		rdsconst.GetBattleSpacePlayerDataKey(uid),
		fmt.Sprintf("%s&%s&%s&%s&%s&%s", uid, playerResults[0].(string), playerResults[2].(string), rdsconst.BattleSpaceStateReady, strconv.FormatInt(int64(index), 10)),
	)

	pipe.Do(ctx, "exec")

	_, err = pipe.Exec(ctx)
	if err != nil {
		pipe.Discard()
		space_mutex.Unlock()
		player_mutex.Unlock()
		return err
	}

	space_mutex.Unlock()
	player_mutex.Unlock()

	return nil
}

func ReadyBattleSpace(ctx context.Context, spaceid string, ready string, clientId *network.ClientID) error {
	uid, err := client.Get(ctx, rdsconst.GetPlayerClientIDKey(clientId.ToString())).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	player_mutex := sync.NewMutex(rdsconst.GetPlayerLockKey(uid))
	if err := player_mutex.Lock(); err != nil {
		return err
	}
	playerResults, err := client.HMGet(ctx, rdsconst.GetPlayerOnlineDataKey(uid),
		rdsconst.PlayerMapClientBattleSpaceId).Result()
	if err != nil {
		player_mutex.Unlock()
		return err
	}
	if len(playerResults) != 1 {
		player_mutex.Unlock()
		return errs.ErrorPlayerOnlineDataLost
	}
	if playerResults[0] != spaceid {
		player_mutex.Unlock()
		return errs.ErrorPlayerIsNotInBattleSpace
	}
	player_mutex.Unlock()

	space_mutex := sync.NewMutex(rdsconst.GetBattleSpaceOnlineDataKey(spaceid))
	if err := space_mutex.Lock(); err != nil {
		return err
	}

	results, err := client.HMGet(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid), rdsconst.GetBattleSpacePlayerDataKey(uid)).Result()
	if err != nil {
		space_mutex.Unlock()
		return err
	}
	if len(results) != 1 {
		return errs.ErrorPlayerIsNotInBattleSpace
	}
	player_data := rdsconst.SplitData(results[0].(string))
	if player_data[3] == ready {
		return errs.ErrorPlayerRepeatOperation
	}
	player_data[3] = ready
	// player_data := makeBattleSpacePlayer(id,icon,display,read,pos)
	pipe := client.TxPipeline()
	defer pipe.Close()

	pipe.Do(ctx, "MULTI")
	pipe.HSet(ctx,
		rdsconst.GetBattleSpaceOnlineDataKey(spaceid),
		rdsconst.GetBattleSpacePlayerDataKey(uid),
		fmt.Sprintf("%s&%s&%s&%s&%s&%s", player_data[0], player_data[1], player_data[2], player_data[3], player_data[4]),
	)

	pipe.Do(ctx, "exec")

	_, err = pipe.Exec(ctx)
	if err != nil {
		pipe.Discard()
		space_mutex.Unlock()
		return err
	}
	space_mutex.Unlock()
	return nil
}

func LeaveBattleSpace(ctx context.Context, clientId *network.ClientID) error {
	uid, err := client.Get(ctx, rdsconst.GetPlayerClientIDKey(clientId.ToString())).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	player_mutex := sync.NewMutex(rdsconst.GetPlayerLockKey(uid))
	if err := player_mutex.Lock(); err != nil {
		return err
	}
	playerResults, err := client.HMGet(ctx, rdsconst.GetPlayerOnlineDataKey(uid),
		rdsconst.PlayerMapClientBattleSpaceId).Result()
	if err != nil {
		player_mutex.Unlock()
		return err
	}
	if len(playerResults) != 1 {
		player_mutex.Unlock()
		return errs.ErrorPlayerOnlineDataLost
	}
	if playerResults[0] == "" {
		player_mutex.Unlock()
		return errs.ErrorPlayerIsNotInBattleSpace
	}
	player_mutex.Unlock()
	spaceid := playerResults[0].(string)
	space_mutex := sync.NewMutex(rdsconst.GetBattleSpaceOnlineDataKey(spaceid))
	if err := space_mutex.Lock(); err != nil {
		return err
	}

	// results, err := client.HGetAll(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid)).Result()
	// if err != nil {
	// 	space_mutex.Unlock()
	// 	return err
	// }

	pipe := client.TxPipeline()
	defer pipe.Close()

	pipe.Do(ctx, "MULTI")

	pipe.SRem(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid), rdsconst.GetBattleSpacePlayerDataKey(uid))
	pipe.HSet(ctx, rdsconst.GetPlayerOnlineDataKey(uid), rdsconst.PlayerMapClientBattleSpaceId, "")

	pipe.Do(ctx, "exec")

	_, err = pipe.Exec(ctx)
	if err != nil {
		pipe.Discard()
		space_mutex.Unlock()
		return err
	}
	space_mutex.Unlock()
	return nil
}

func GetBattleSpaceList(ctx context.Context, start int64, size int64) ([]string, error) {

	val, err := client.LRange(ctx, rdsconst.GetBattleSpaceOnlineDataKey(""), start, (start+1)*size).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func AutoEnterBattleSpace(ctx context.Context, clientId *network.ClientID) error {
	return nil
}

func GetBattleSpacesCount(ctx context.Context) (int64, error) {
	return client.LLen(ctx, rdsconst.GetBattleSpaceOnlineDataKey("")).Result()
}

func FindBattleSpaceData(ctx context.Context, spaceid string) (map[string]string, error) {
	mutex := sync.NewMutex(rdsconst.GetBattleSpaceOnlineDataKey(spaceid))
	if err := mutex.Lock(); err != nil {
		return nil, err
	}
	results, err := client.HGetAll(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid)).Result()
	if err != nil {
		mutex.Unlock()
		return nil, err
	}
	mutex.Unlock()

	return results, nil
}

func IsMaster(ctx context.Context, clientId *network.ClientID) (bool, error) {

	uid, err := client.Get(ctx, rdsconst.GetPlayerClientIDKey(clientId.ToString())).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	player_mutex := sync.NewMutex(rdsconst.GetPlayerLockKey(uid))
	if err := player_mutex.Lock(); err != nil {
		return false, err
	}
	playerResults, err := client.HMGet(ctx, rdsconst.GetPlayerOnlineDataKey(uid),
		rdsconst.PlayerMapClientBattleSpaceId).Result()
	if err != nil {
		player_mutex.Unlock()
		return false, err
	}
	if len(playerResults) != 1 {
		player_mutex.Unlock()
		return false, errs.ErrorPlayerOnlineDataLost
	}
	spaceid := playerResults[0].(string)

	mutex := sync.NewMutex(rdsconst.GetBattleSpaceOnlineDataKey(spaceid))
	if err := mutex.Lock(); err != nil {
		player_mutex.Unlock()
		return false, err
	}
	results, err := client.HGetAll(ctx, rdsconst.GetBattleSpaceOnlineDataKey(spaceid)).Result()
	if err != nil {
		player_mutex.Unlock()
		mutex.Unlock()
		return false, err
	}
	player_mutex.Unlock()
	mutex.Unlock()
	return clientId.Equal(&network.ClientID{Address: results[rdsconst.PalyerMapClientIdAddress],
		Id: results[rdsconst.PlayerMapClientIdId]}), nil
}

func createBattleSpacePlayerPos(list []string) string {
	res := ""
	for i := (0); i < len(list); i++ {
		res += fmt.Sprintf("%s&", list[i])
	}
	return res[:len(res)-1]
}

func updateBattleSpacePlayerPos(poss string, uid string, pos int) string {
	list := strings.Split(poss, "&")
	for i := 0; i < len(list); i++ {
		if i == pos {
			list[i] = uid
		}
	}
	return createBattleSpacePlayerPos(list)
}

func enterBattleSpacePlayerPos(poss string, uid string) (string, int32) {
	list := strings.Split(poss, "&")
	pos := int32(0)
	for i := 0; i < len(list); i++ {
		if list[i] == "" {
			list[i] = uid
			pos = int32(i)
			break
		}
	}
	return createBattleSpacePlayerPos(list), pos
}
