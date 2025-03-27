import React, { useState } from 'react';
import axios from 'axios';
import {
    Box,
    Container,
    Paper,
    TextField,
    Button,
    Typography,
    Alert,
    Link,
    Divider,
    CircularProgress,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { styled } from '@mui/material/styles';

const API_URL = process.env.REACT_APP_API_URL || 'https://ordermanagersystem-auth-service.onrender.com';

const StyledContainer = styled(Container)(({ theme }) => ({
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: 'linear-gradient(135deg, #6B73FF 0%, #000DFF 100%)',
    padding: theme.spacing(3),
}));

const StyledPaper = styled(Paper)(({ theme }) => ({
    padding: theme.spacing(4),
    width: '100%',
    maxWidth: '400px',
    borderRadius: '16px',
    boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
}));

const StyledButton = styled(Button)(({ theme }) => ({
    padding: theme.spacing(1.5),
    borderRadius: '8px',
    textTransform: 'none',
    fontSize: '1rem',
    fontWeight: 600,
    boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
    '&:hover': {
        boxShadow: '0 6px 16px rgba(0, 0, 0, 0.2)',
    },
}));

const StyledTextField = styled(TextField)(({ theme }) => ({
    '& .MuiOutlinedInput-root': {
        borderRadius: '8px',
        '&:hover fieldset': {
            borderColor: theme.palette.primary.main,
        },
    },
}));

const LoadingContainer = styled(Box)({
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: 'rgba(255, 255, 255, 0.9)',
    backdropFilter: 'blur(4px)',
    zIndex: 999,
    borderRadius: '16px',
});

const StyledCircularProgress = styled(CircularProgress)({
    color: '#4C6EF5',
    size: 50,
});

// 創建一個 axios 實例
const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
        'Accept': '*/*'
    },
    withCredentials: false
});

const Login = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        email: 'admin@example.com',
        password: 'password123',
    });
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    const handleLogin = async (e) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            const response = await api.post('/api/v1/auth/login', formData);
            const { token, user } = response.data;

            if (token) {
                localStorage.setItem('userToken', token);
                localStorage.setItem('userData', JSON.stringify(user));
                // 觸發登入狀態變更事件
                window.dispatchEvent(new Event('loginStateChange'));
                const redirectUrl = sessionStorage.getItem('redirectUrl');
                sessionStorage.removeItem('redirectUrl');
                navigate(redirectUrl || '/');
            } else {
                throw new Error('登入回應中沒有 token');
            }
        } catch (err) {
            console.error('Login error:', err);
            setError(
                err.response?.data?.message ||
                    err.message === 'Network Error'
                    ? '無法連接到伺服器，請檢查您的網路連接'
                    : '登入失敗，請檢查您的帳號密碼是否正確'
            );
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <StyledContainer maxWidth={false} disableGutters>
            <StyledPaper elevation={3}>
                {isLoading && (
                    <LoadingContainer>
                        <StyledCircularProgress size={50} thickness={4} />
                    </LoadingContainer>
                )}
                <Box sx={{ textAlign: 'center', mb: 4 }}>
                    <Typography variant="h4" component="h1" sx={{ fontWeight: 700, color: 'primary.main' }}>
                        訂單管理系統
                    </Typography>
                    <Typography variant="subtitle1" sx={{ mt: 1, color: 'text.secondary' }}>
                        歡迎回來！請登入您的帳號
                    </Typography>
                </Box>

                {error && (
                    <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>
                        {error}
                    </Alert>
                )}

                <Box component="form" onSubmit={handleLogin}>
                    <StyledTextField
                        margin="normal"
                        required
                        fullWidth
                        id="email"
                        label="電子郵件"
                        name="email"
                        autoComplete="email"
                        autoFocus
                        value={formData.email}
                        onChange={handleChange}
                    />
                    <StyledTextField
                        margin="normal"
                        required
                        fullWidth
                        name="password"
                        label="密碼"
                        type="password"
                        id="password"
                        autoComplete="current-password"
                        value={formData.password}
                        onChange={handleChange}
                    />
                    <StyledButton
                        type="submit"
                        fullWidth
                        variant="contained"
                        sx={{ mt: 3 }}
                        disabled={isLoading}
                    >
                        {isLoading ? '登入中...' : '登入'}
                    </StyledButton>

                    <Divider sx={{ my: 3 }}>
                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                            或
                        </Typography>
                    </Divider>

                    <Box sx={{ textAlign: 'center' }}>
                        <Link
                            component="button"
                            variant="body2"
                            onClick={() => navigate('/register')}
                            sx={{
                                mb: 1,
                                color: 'primary.main',
                                textDecoration: 'none',
                                '&:hover': {
                                    textDecoration: 'underline',
                                }
                            }}
                        >
                            還沒有帳號？立即註冊
                        </Link>
                        <br />
                        <Link
                            component="button"
                            variant="body2"
                            onClick={() => navigate('/forgot-password')}
                            sx={{
                                color: 'text.secondary',
                                textDecoration: 'none',
                                '&:hover': {
                                    textDecoration: 'underline',
                                }
                            }}
                        >
                            忘記密碼？
                        </Link>
                    </Box>
                </Box>
            </StyledPaper>
        </StyledContainer>
    );
};

export default Login;