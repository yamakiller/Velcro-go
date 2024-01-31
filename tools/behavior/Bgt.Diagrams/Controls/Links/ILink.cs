﻿using System.Windows;

namespace Bgt.Diagrams.Controls
{
    public interface ILink
    {
        IPort Source { get; set; }
        IPort Target { get; set; }

        IPort Control1 { get; set; }
        IPort Control2 { get; set; }

        Point SourcePoint { get; set; }
        Point TargetPoint { get; set; }

        Point ControlPoint1 { get; set; }
        Point ControlPoint2 { get; set; }

        void UpdatePath();
    }
}
