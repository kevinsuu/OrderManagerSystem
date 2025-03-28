import React from 'react';
import { HashRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme';
import { useHealthCheck } from './utils/healthCheck';
import { Snackbar, Alert } from '@mui/material';

// Layouts
import StoreLayout from './layouts/Store/StoreLayout';
import AdminLayout from './layouts/Admin/AdminLayout';

// Store Pages
import Home from './pages/Store/Home';
import Login from './pages/Store/Login';
import Register from './pages/Store/Register';
import ForgotPassword from './pages/Store/ForgotPassword';
import Profile from './pages/Store/Profile';
import Cart from './pages/Store/Cart';
import ProductDetail from './pages/Store/ProductDetail';

// Admin Pages
import AdminLogin from './pages/Admin/Login';
import AdminDashboard from './pages/Admin/Dashboard';

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
    const { isAnyServiceDown } = useHealthCheck();

    return (
        <>
            <ThemeProvider theme={theme}>
                <CssBaseline />
                <Router>
                    <Routes>
                        {/* 登入相關頁面 */}
                        <Route path="/login" element={<Login />} />
                        <Route path="/register" element={<Register />} />
                        <Route path="/forgot-password" element={<ForgotPassword />} />

                        {/* 商店前台路由 */}
                        <Route path="/" element={<StoreLayout />}>
                            <Route index element={<Home />} />
                            <Route
                                path="profile"
                                element={
                                    <UserProtectedRoute>
                                        <Profile />
                                    </UserProtectedRoute>
                                }
                            />
                            <Route path="store" element={<Home />} />
                            <Route path="store/product/:id" element={<ProductDetail />} />
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
                            <Route path="cart" element={<Cart />} />
                        </Route>

                        {/* 管理後台路由 */}
                        <Route path="/admin" element={<AdminLayout />}>
                            <Route path="login" element={<AdminLogin />} />
                            <Route
                                path="dashboard"
                                element={
                                    <ProtectedRoute>
                                        <AdminDashboard />
                                    </ProtectedRoute>
                                }
                            />
                        </Route>
                    </Routes>
                </Router>

                {/* 服務狀態提示 */}
                <Snackbar
                    open={isAnyServiceDown}
                    anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                    sx={{ position: 'fixed', bottom: 16, right: 16 }}
                >
                    <Alert severity="error">
                        系統服務異常，請稍後再試
                    </Alert>
                </Snackbar>
            </ThemeProvider>
        </>
    );
}

export default App; 