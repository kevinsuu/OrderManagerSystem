import React, { useState, useEffect } from 'react';
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
    InputAdornment,
    IconButton,
} from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import { styled } from '@mui/material/styles';
import { Link as RouterLink } from 'react-router-dom';
import { Visibility, VisibilityOff } from '@mui/icons-material';

const API_URL = process.env.REACT_APP_AUTH_SERVICE_URL;

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

const LoadingOverlay = styled(Box)({
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    background: 'rgba(0, 13, 255, 0.05)',
    backdropFilter: 'blur(8px)',
    zIndex: 999,
    borderRadius: '16px',
    overflow: 'hidden',
});

const LoadingText = styled(Typography)({
    color: '#4C6EF5',
    marginTop: '16px',
    fontSize: '0.875rem',
    fontWeight: 500,
    zIndex: 2,
});

const PulseContainer = styled(Box)({
    position: 'relative',
    width: '80px',
    height: '80px',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
});

const PulseRing = styled(Box)(({ delay = 0 }) => ({
    position: 'absolute',
    width: '100%',
    height: '100%',
    border: '3px solid #4C6EF5',
    borderRadius: '50%',
    animation: 'pulse 2s cubic-bezier(0.455, 0.03, 0.515, 0.955) infinite',
    animationDelay: `${delay}ms`,
    opacity: 0,
    '@keyframes pulse': {
        '0%': {
            transform: 'scale(0.4)',
            opacity: 0,
        },
        '50%': {
            opacity: 0.5,
        },
        '100%': {
            transform: 'scale(1.2)',
            opacity: 0,
        },
    },
}));

const InnerCircle = styled(Box)({
    width: '20px',
    height: '20px',
    backgroundColor: '#4C6EF5',
    borderRadius: '50%',
    animation: 'glow 1.5s ease-in-out infinite',
    '@keyframes glow': {
        '0%, 100%': {
            transform: 'scale(1)',
            opacity: 1,
        },
        '50%': {
            transform: 'scale(1.2)',
            opacity: 0.7,
        },
    },
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
    const location = useLocation();
    const [formData, setFormData] = useState({
        email: 'admin@example.com',
        password: 'password123',
    });
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    useEffect(() => {
        if (location.state?.message) {
            if (location.state.severity === 'success') {
                setSuccess(location.state.message);
            } else {
                setError(location.state.message);
            }
            navigate(location.pathname, { replace: true });
        }
    }, [location, navigate]);

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
                    <LoadingOverlay>
                        <PulseContainer>
                            <PulseRing delay={0} />
                            <PulseRing delay={400} />
                            <PulseRing delay={800} />
                            <InnerCircle />
                        </PulseContainer>
                        <LoadingText>處理中...</LoadingText>
                    </LoadingOverlay>
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

                {success && (
                    <Alert severity="success" sx={{ mb: 3, borderRadius: 2 }}>
                        {success}
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
                        type={showPassword ? 'text' : 'password'}
                        id="password"
                        autoComplete="current-password"
                        value={formData.password}
                        onChange={handleChange}
                        InputProps={{
                            endAdornment: (
                                <InputAdornment position="end">
                                    <IconButton
                                        onClick={() => setShowPassword(!showPassword)}
                                        edge="end"
                                    >
                                        {showPassword ? <VisibilityOff /> : <Visibility />}
                                    </IconButton>
                                </InputAdornment>
                            ),
                        }}
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
                            component={RouterLink}
                            to="/register"
                            variant="body2"
                            sx={{
                                mb: 1,
                                color: 'primary.main',
                                textDecoration: 'none',
                                display: 'block',
                                '&:hover': {
                                    textDecoration: 'underline',
                                }
                            }}
                        >
                            還沒有帳號？立即註冊
                        </Link>
                        <Link
                            component={RouterLink}
                            to="/forgot-password"
                            variant="body2"
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