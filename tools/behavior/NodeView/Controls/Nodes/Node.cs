﻿using Bga.Diagrams.Adorners;
using Bga.Diagrams.Controls.Ports;
using Bga.Diagrams.Controls;

using System.Windows.Controls;
using System.Windows.Documents;
using System.Windows;

namespace Bga.Diagrams.Controls
{
    public class Node : DiagramItem, INode
    {
        static Node()
        {
            FrameworkElement.DefaultStyleKeyProperty.OverrideMetadata(
                typeof(Node), new FrameworkPropertyMetadata(typeof(Node)));
        }

        #region Properties

        #region Content Property

        public object Content
        {
            get { return (bool)GetValue(ContentProperty); }
            set { SetValue(ContentProperty, value); }
        }

        public static readonly DependencyProperty ContentProperty =
            DependencyProperty.Register("Content",
                                       typeof(object),
                                       typeof(Node));

        #endregion

        #region CanResize Property

        public bool CanResize
        {
            get { return (bool)GetValue(CanResizeProperty); }
            set { SetValue(CanResizeProperty, value); }
        }

        public static readonly DependencyProperty CanResizeProperty =
            DependencyProperty.Register("CanResize",
                                       typeof(bool),
                                       typeof(Node),
                                       new FrameworkPropertyMetadata(true));

        #endregion

        private List<IPort> m_ports = new List<IPort>();
        public ICollection<IPort> Ports { get { return m_ports; } }

        public override Rect Bounds
        {
            get
            {
                //var itemRect = VisualTreeHelper.GetDescendantBounds(item);
                //return item.TransformToAncestor(this).TransformBounds(itemRect);
                var x = Canvas.GetLeft(this);
                var y = Canvas.GetTop(this);
                return new Rect(x, y, ActualWidth, ActualHeight);
            }
        }

        #endregion

        public Node()
        {
        }

        public void UpdatePosition()
        {
            foreach (var p in Ports)
                p.UpdatePosition();
        }

        protected override Adorner CreateSelectionAdorner()
        {
            return new SelectionAdorner(this, new SelectionFrame());
        }

        #region INode Members

        IEnumerable<IPort> INode.Ports
        {
            get { return m_ports; }
        }

        #endregion
    }
}
