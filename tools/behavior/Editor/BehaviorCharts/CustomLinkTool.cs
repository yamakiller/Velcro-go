﻿using Bgt.Diagrams.Controls;
using Bgt.Diagrams.Tools;
using Bgt.Diagrams;
using System.Windows;

namespace Editor.BehaviorCharts
{
    class CustomLinkTool : LinkTool
    {
        public CustomLinkTool(DiagramView view)
            : base(view)
        {
        }

        protected override ILink CreateNewLink(IPort port)
        {
            var link = new OrthogonalLink();
            BindNewLinkToPort(port, link);
            return link;
        }

        protected override void UpdateLink(Point point, IPort port)
        {
            base.UpdateLink(point, port);
            var link = Link as OrthogonalLink;
        }
    }
}
