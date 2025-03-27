import React, { useState } from 'react';
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
    InputAdornment,
    IconButton
} from '@mui/material';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import { Visibility, VisibilityOff } from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import axios from 'axios';

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

const API_URL = process.env.REACT_APP_AUTH_SERVICE_URL || 'https://ordermanagersystem-auth-service.onrender.com';

// 創建一個 axios 實例
const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
        'Accept': '*/*'
    },
    withCredentials: false
});

const Register = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        confirmPassword: ''
    });
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
        setError('');
    };

    const validateForm = () => {
        if (!formData.username || !formData.email || !formData.password || !formData.confirmPassword) {
            setError('請填寫所有欄位');
            return false;
        }
        if (formData.password !== formData.confirmPassword) {
            setError('密碼與確認密碼不符');
            return false;
        }
        if (formData.password.length < 3) {
            setError('密碼長度至少需要3個字元');
            return false;
        }
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
            setError('請輸入有效的電子郵件地址');
            return false;
        }
        return true;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!validateForm()) return;

        setLoading(true);
        try {
            const response = await api.post('/api/v1/auth/register', {
                username: formData.username,
                email: formData.email,
                password: formData.password
            });

            if (response.data) {
                navigate('/login', {
                    state: {
                        message: '註冊成功！請使用您的帳號密碼登入',
                        severity: 'success'
                    }
                });
            }
        } catch (error) {
            // 檢查錯誤回應中的 error 欄位
            if (error.response?.data?.error === 'username already taken') {
                setError('此使用者名稱已被使用，請選擇其他名稱');
            } else if (error.response?.data?.error === 'user already exists') {
                setError('此電子郵件已被註冊，請使用其他電子郵件');
            } else if (error.message === 'Network Error') {
                setError('無法連接到伺服器，請檢查您的網路連接');
            } else {
                setError('註冊失敗，請稍後再試');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <StyledContainer maxWidth={false} disableGutters>
            <StyledPaper elevation={3}>
                {loading && (
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
                        創建新帳號
                    </Typography>
                </Box>

                {error && (
                    <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>
                        {error}
                    </Alert>
                )}

                <Box component="form" onSubmit={handleSubmit}>
                    <StyledTextField
                        margin="normal"
                        required
                        fullWidth
                        id="username"
                        label="使用者名稱"
                        name="username"
                        autoComplete="username"
                        autoFocus
                        value={formData.username}
                        onChange={handleChange}
                    />
                    <StyledTextField
                        margin="normal"
                        required
                        fullWidth
                        id="email"
                        label="電子郵件"
                        name="email"
                        autoComplete="email"
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
                        autoComplete="new-password"
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
                    <StyledTextField
                        margin="normal"
                        required
                        fullWidth
                        name="confirmPassword"
                        label="確認密碼"
                        type={showConfirmPassword ? 'text' : 'password'}
                        id="confirmPassword"
                        autoComplete="new-password"
                        value={formData.confirmPassword}
                        onChange={handleChange}
                        InputProps={{
                            endAdornment: (
                                <InputAdornment position="end">
                                    <IconButton
                                        onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                        edge="end"
                                    >
                                        {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
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
                        disabled={loading}
                    >
                        {loading ? '註冊中...' : '註冊'}
                    </StyledButton>

                    <Divider sx={{ my: 3 }}>
                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                            或
                        </Typography>
                    </Divider>

                    <Box sx={{ textAlign: 'center' }}>
                        <Link
                            component={RouterLink}
                            to="/login"
                            variant="body2"
                            sx={{
                                color: 'primary.main',
                                textDecoration: 'none',
                                '&:hover': {
                                    textDecoration: 'underline',
                                }
                            }}
                        >
                            已有帳號？立即登入
                        </Link>
                    </Box>
                </Box>
            </StyledPaper>
        </StyledContainer>
    );
};

export default Register; 