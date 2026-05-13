import { useState, useEffect, useRef } from 'react';
import { getNotifications, markAsRead, markAllAsRead } from '../../services/notificationService';
import { useAuth } from '../../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import './NotificationBell.css';

const NotificationBell = () => {
  const { isLoggedIn } = useAuth();
  const navigate = useNavigate();
  const [isOpen, setIsOpen] = useState(false);
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const dropdownRef = useRef(null);

  const fetchNotifications = async () => {
    if (!isLoggedIn) return;
    try {
      setLoading(true);
      const res = await getNotifications({ limit: 10 });
      if (res.success) {
        setNotifications(res.data.notifications || []);
        setUnreadCount(res.data.unread_count || 0);
      }
    } catch (err) {
      console.error('Lỗi tải thông báo:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNotifications();
    // Setup polling every 3 minutes
    const interval = setInterval(fetchNotifications, 3 * 60 * 1000);
    return () => clearInterval(interval);
  }, [isLoggedIn]);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleToggle = () => {
    if (!isOpen && isLoggedIn) {
      fetchNotifications();
    }
    setIsOpen(!isOpen);
  };

  const handleMarkAsRead = async (id, e) => {
    e.stopPropagation();
    try {
      await markAsRead(id);
      setNotifications(prev => prev.map(n => n.id === id ? { ...n, is_read: true } : n));
      setUnreadCount(prev => Math.max(0, prev - 1));
    } catch (err) {
      console.error(err);
    }
  };

  const handleMarkAllAsRead = async () => {
    try {
      await markAllAsRead();
      setNotifications(prev => prev.map(n => ({ ...n, is_read: true })));
      setUnreadCount(0);
    } catch (err) {
      console.error(err);
    }
  };

  const handleNotifClick = (notif) => {
    if (!notif.is_read) {
      handleMarkAsRead(notif.id, { stopPropagation: () => {} });
    }
    setIsOpen(false);
    
    // Điều hướng dựa trên loại thông báo
    if (notif.type === 'booking' || notif.type === 'payment') {
      navigate('/account/bookings');
    } else if (notif.type === 'review') {
      // Navigate to tours? or bookings?
      // Since review is linked to a tour, maybe they can just see it in their bookings.
      navigate('/account/bookings');
    }
  };

  if (!isLoggedIn) return null;

  return (
    <div className="notif-wrapper" ref={dropdownRef}>
      <button className="notif-bell" onClick={handleToggle}>
        <span className="notif-icon">🔔</span>
        {unreadCount > 0 && (
          <span className="notif-badge">{unreadCount > 99 ? '99+' : unreadCount}</span>
        )}
      </button>

      {isOpen && (
        <div className="notif-dropdown">
          <div className="notif-header">
            <h4>Thông báo</h4>
            {unreadCount > 0 && (
              <button className="mark-all-btn" onClick={handleMarkAllAsRead}>
                Đánh dấu tất cả đã đọc
              </button>
            )}
          </div>
          
          <div className="notif-list">
            {loading && notifications.length === 0 ? (
              <div className="notif-empty">Đang tải...</div>
            ) : notifications.length === 0 ? (
              <div className="notif-empty">Bạn không có thông báo nào</div>
            ) : (
              notifications.map(notif => (
                <div 
                  key={notif.id} 
                  className={`notif-item ${!notif.is_read ? 'unread' : ''}`}
                  onClick={() => handleNotifClick(notif)}
                >
                  <div className="notif-icon-circle">
                    {notif.type === 'booking' ? '📅' : notif.type === 'payment' ? '💳' : notif.type === 'review' ? '⭐' : '👋'}
                  </div>
                  <div className="notif-content">
                    <h5>{notif.title}</h5>
                    <p>{notif.message}</p>
                    <span className="notif-time">
                      {new Date(notif.created_at).toLocaleString('vi-VN')}
                    </span>
                  </div>
                  {!notif.is_read && (
                    <div 
                      className="notif-dot" 
                      title="Đánh dấu đã đọc"
                      onClick={(e) => handleMarkAsRead(notif.id, e)}
                    />
                  )}
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default NotificationBell;
