'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function Notifications() {
  const [notifications, setNotifications] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchNotifications();
  }, []);

  const fetchNotifications = async () => {
    try {
      const response = await apiClients.notification.get('/api/v1/notifications');
      setNotifications(response.data || []);
    } catch (error) {
      console.error('Failed to fetch notifications:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Notifications</h1>

      {loading ? (
        <div>Loading...</div>
      ) : notifications.length === 0 ? (
        <Card>
          <p className="text-gray-600">No notifications.</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {notifications.map((notif) => (
            <Card key={notif.id} className={!notif.is_read ? 'bg-blue-50' : ''}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{notif.title}</h3>
                  <p className="text-gray-600 text-sm mt-1">{notif.message}</p>
                  <p className="text-xs text-gray-500 mt-2">
                    {new Date(notif.created_at).toLocaleString()}
                  </p>
                </div>
                {!notif.is_read && (
                  <span className="bg-blue-600 text-white text-xs px-2 py-1 rounded-full">New</span>
                )}
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
