'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function Chat() {
  const [threads, setThreads] = useState<any[]>([]);
  const [selectedThread, setSelectedThread] = useState<string | null>(null);
  const [messages, setMessages] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchThreads();
  }, []);

  useEffect(() => {
    if (selectedThread) {
      fetchMessages(selectedThread);
    }
  }, [selectedThread]);

  const fetchThreads = async () => {
    try {
      const response = await apiClients.collaboration.get('/api/v1/threads/user');
      setThreads(response.data || []);
    } catch (error) {
      console.error('Failed to fetch threads:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchMessages = async (threadId: string) => {
    try {
      const response = await apiClients.collaboration.get(`/api/v1/threads/${threadId}/messages`);
      setMessages(response.data || []);
    } catch (error) {
      console.error('Failed to fetch messages:', error);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Chat</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-1">
          <Card title="Conversations">
            {loading ? (
              <div>Loading...</div>
            ) : threads.length === 0 ? (
              <p className="text-gray-600">No conversations yet.</p>
            ) : (
              <div className="space-y-2">
                {threads.map((thread) => (
                  <button
                    key={thread.id}
                    onClick={() => setSelectedThread(thread.id)}
                    className={`w-full text-left p-3 rounded-lg hover:bg-gray-50 ${selectedThread === thread.id ? 'bg-primary-50' : ''}`}
                  >
                    <div className="font-semibold">{thread.title || 'Untitled'}</div>
                    <div className="text-sm text-gray-500">Last message...</div>
                  </button>
                ))}
              </div>
            )}
          </Card>
        </div>

        <div className="lg:col-span-2">
          <Card title={selectedThread ? 'Messages' : 'Select a conversation'}>
            {selectedThread ? (
              <div className="space-y-4">
                {messages.map((msg) => (
                  <div key={msg.id} className="p-3 bg-gray-50 rounded-lg">
                    <div className="font-semibold text-sm">{msg.sender_id}</div>
                    <div className="text-gray-700">{msg.content}</div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-600">Select a conversation to view messages</p>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
}
