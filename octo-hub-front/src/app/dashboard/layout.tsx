'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { 
  HomeIcon, 
  UsersIcon, 
  ChevronLeftIcon,
  ChevronRightIcon,
  Bars3Icon,
  BellIcon,
  UserCircleIcon,
  WifiIcon
} from '@heroicons/react/24/outline';
import { apiGet, TokenManager, ApiResponse } from '@/utils/api';
import { useWebSocket, getWebSocketStatusText, getWebSocketStatusColor } from '@/hooks/useWebSocket';

interface UserInfo {
  id: number;
  email: string;
}

const menuItems = [
  { name: '仪表盘', href: '/dashboard', icon: HomeIcon },
  { name: '账号列表', href: '/dashboard/accounts', icon: UsersIcon },
];

export default function DashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const router = useRouter();
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null);
  const [loading, setLoading] = useState(true);

  // WebSocket连接管理
  const { 
    status: wsStatus, 
    isConnected: wsConnected, 
    connect: wsConnect, 
    disconnect: wsDisconnect,
    sendMessage,
    lastMessage,
    connectionCount
  } = useWebSocket({
    autoConnect: false, // 等待用户信息加载完成后再连接
    pingInterval: 10000, // 50秒心跳 - 比服务端60秒超时短
    reconnectInterval: 5000, // 5秒重连间隔
    maxReconnectAttempts: 10 // 最大重连10次
  });

  // 获取用户信息
  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const response = await apiGet<ApiResponse>('/user/info');
        setUserInfo(response.data);
      } catch (error) {
        console.error('获取用户信息失败:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchUserInfo();
  }, [router]);

  // 用户信息加载完成后自动连接WebSocket
  useEffect(() => {
    if (userInfo && !loading) {
      console.log('User info loaded, connecting to WebSocket...');
      wsConnect().catch(error => {
        console.error('Failed to connect to WebSocket:', error);
      });
    }
  }, [userInfo, loading, wsConnect]);

  // 处理WebSocket消息
  useEffect(() => {
    if (lastMessage) {
      console.log('Received WebSocket message in Dashboard:', lastMessage);
      // 这里可以根据消息类型进行不同的处理
      // 例如：显示通知、更新状态等
    }
  }, [lastMessage]);

  // 退出登录
  const handleLogout = () => {
    // 断开WebSocket连接
    wsDisconnect();
    TokenManager.removeToken();
    router.push('/login');
  };

  // 手动重连WebSocket
  const handleWebSocketReconnect = () => {
    wsConnect().catch(error => {
      console.error('Manual reconnect failed:', error);
    });
  };

  // 加载中状态
  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-100">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Mobile menu overlay */}
      {mobileMenuOpen && (
        <div 
          className="fixed inset-0 z-40 bg-gray-600 bg-opacity-75 lg:hidden"
          onClick={() => setMobileMenuOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div className={`
        fixed inset-y-0 left-0 z-50 bg-white shadow-lg transform transition-all duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0
        ${mobileMenuOpen ? 'translate-x-0' : '-translate-x-full'}
        ${sidebarCollapsed ? 'w-16' : 'w-64'}
        lg:flex lg:flex-col
      `}>
        {/* Sidebar header */}
        <div className="flex items-center justify-between h-16 px-4 border-b border-gray-200">
          {!sidebarCollapsed && (
            <div className="flex items-center">
              <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold">O</span>
              </div>
              <span className="ml-2 text-xl font-semibold text-gray-800">OctoHub</span>
            </div>
          )}
          <button
            onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
            className="hidden lg:flex p-2 rounded-md text-gray-500 hover:text-gray-700 hover:bg-gray-100"
          >
            {sidebarCollapsed ? (
              <ChevronRightIcon className="w-5 h-5" />
            ) : (
              <ChevronLeftIcon className="w-5 h-5" />
            )}
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-4 py-4 space-y-2">
          {menuItems.map((item) => (
            <Link
              key={item.name}
              href={item.href}
              className="flex items-center px-3 py-2 text-sm font-medium text-gray-700 rounded-md hover:text-gray-900 hover:bg-gray-100 transition-colors"
            >
              <item.icon className={`${sidebarCollapsed ? 'w-6 h-6' : 'w-5 h-5 mr-3'} flex-shrink-0`} />
              {!sidebarCollapsed && <span>{item.name}</span>}
            </Link>
          ))}
        </nav>

        {/* Sidebar footer */}
        {!sidebarCollapsed && userInfo && (
          <div className="p-4 border-t border-gray-200">
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <UserCircleIcon className="w-8 h-8 text-gray-400" />
                <div className="ml-3">
                  <p className="text-sm font-medium text-gray-700">
                    {`用户${userInfo.id}`}
                  </p>
                  <p className="text-xs text-gray-500">{userInfo.email}</p>
                </div>
              </div>
              <button
                onClick={handleLogout}
                className="text-xs text-gray-500 hover:text-red-600 transition-colors"
                title="退出登录"
              >
                退出
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Main content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top header */}
        <header className="bg-white shadow-sm border-b border-gray-200">
          <div className="flex items-center justify-between h-16 px-4">
            <div className="flex items-center">
              <button
                onClick={() => setMobileMenuOpen(true)}
                className="p-2 rounded-md text-gray-500 hover:text-gray-700 hover:bg-gray-100 lg:hidden"
              >
                <Bars3Icon className="w-6 h-6" />
              </button>
            </div>
            
            <div className="flex items-center space-x-4">
              {/* WebSocket连接状态指示器 */}
              <div className="flex items-center space-x-2">
                <button
                  onClick={handleWebSocketReconnect}
                  className={`p-2 rounded-md transition-colors ${
                    wsConnected 
                      ? 'text-green-500 hover:text-green-600 hover:bg-green-50' 
                      : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'
                  }`}
                  title={`WebSocket: ${getWebSocketStatusText(wsStatus)} (点击重连)`}
                >
                  <WifiIcon className="w-5 h-5" />
                </button>
                <span className={`text-xs font-medium hidden sm:block ${getWebSocketStatusColor(wsStatus)}`}>
                  {getWebSocketStatusText(wsStatus)}
                </span>
                {connectionCount > 0 && (
                  <span className="text-xs text-gray-400 hidden md:block">
                    ({connectionCount})
                  </span>
                )}
              </div>

              <button className="p-2 rounded-md text-gray-500 hover:text-gray-700 hover:bg-gray-100">
                <BellIcon className="w-6 h-6" />
              </button>
              <div className="flex items-center">
                <UserCircleIcon className="w-8 h-8 text-gray-400" />
                {userInfo && (
                  <span className="ml-2 text-sm font-medium text-gray-700 hidden sm:block">
                    {`用户${userInfo.id}`}
                  </span>
                )}
              </div>
            </div>
          </div>
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-x-hidden overflow-y-auto bg-gray-50 p-6">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}
