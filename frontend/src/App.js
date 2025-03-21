import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme';

// Layouts
import MainLayout from './layouts/MainLayout';
import StoreLayout from './layouts/StoreLayout';

// Store Pages
import StorePage from './pages/Store/Home';
import Login from './pages/Store/Login';

// Admin Pages
import AdminLogin from './pages/Admin/Login';
import Dashboard from './pages/Admin/Dashboard';

// 路由保護組件
const ProtectedRoute = ({ children }) => {
    const isAuthenticated = localStorage.getItem('adminLoggedIn') === 'true';

    if (!isAuthenticated) {
        return <Navigate to="/admin/login" />;
    }

    return children;
};

// 使用者路由保護組件
const UserProtectedRoute = ({ children }) => {
    const isAuthenticated = localStorage.getItem('userToken');

    if (!isAuthenticated) {
        // 儲存當前嘗試訪問的路徑，登入後可以重定向回來
        sessionStorage.setItem('redirectUrl', window.location.pathname);
        return <Navigate to="/login" />;
    }

    return children;
};

function App() {
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <Routes>
                {/* 商店前台路由 */}
                <Route path="/" element={<StoreLayout />}>
                    <Route index element={<StorePage />} />
                    <Route path="store" element={<StorePage />} />
                    <Route path="login" element={<Login />} />

                    {/* 需要登入的路由 */}
                    <Route path="store/profile" element={
                        <UserProtectedRoute>
                            <div>個人資料頁面</div>
                        </UserProtectedRoute>
                    } />
                    <Route path="store/orders" element={
                        <UserProtectedRoute>
                            <div>訂單記錄頁面</div>
                        </UserProtectedRoute>
                    } />
                    <Route path="store/wishlist" element={
                        <UserProtectedRoute>
                            <div>收藏清單頁面</div>
                        </UserProtectedRoute>
                    } />
                    <Route path="store/cart" element={
                        <UserProtectedRoute>
                            <div>購物車頁面</div>
                        </UserProtectedRoute>
                    } />
                </Route>

                {/* 管理後台路由 */}
                <Route path="/admin/login" element={<AdminLogin />} />
                <Route
                    path="/admin/*"
                    element={
                        <ProtectedRoute>
                            <MainLayout />
                        </ProtectedRoute>
                    }
                >
                    <Route path="dashboard" element={<Dashboard />} />
                </Route>
            </Routes>
        </ThemeProvider>
    );
}

export default App; 