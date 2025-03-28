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
    IconButton
} from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { styled } from '@mui/material/styles';
import axios from 'axios';
import VisibilityIcon from '@mui/icons-material/Visibility';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';

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

const ForgotPassword = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        email: '',
        newPassword: '',
        confirmPassword: ''
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [showPassword, setShowPassword] = useState({
        newPassword: false,
        confirmPassword: false
    });
    const AUTH_SERVICE_URL = process.env.REACT_APP_AUTH_SERVICE_URL || 'https://ordermanagersystem-auth-service.onrender.com';
    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const { email, newPassword, confirmPassword } = formData;

        // 基本驗證
        if (!email || !newPassword || !confirmPassword) {
            setError('請填寫所有欄位');
            return;
        }
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
            setError('請輸入有效的電子郵件地址');
            return;
        }
        if (newPassword !== confirmPassword) {
            setError('兩次輸入的密碼不一致');
            return;
        }
        if (newPassword.length < 6) {
            setError('密碼長度至少需要6個字符');
            return;
        }

        setLoading(true);
        setError('');
        setSuccess('');

        try {
            await axios.post(`${AUTH_SERVICE_URL}/api/v1/auth/forgot-password`, {
                email,
                newPassword,
                confirmPassword
            });

            setSuccess('密碼重置成功！');
            setTimeout(() => {
                navigate('/login');
            }, 2000);
        } catch (error) {
            setError(error.response?.data?.error || '密碼重置失敗，請稍後再試');
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
                        重設密碼
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

                <Box component="form" onSubmit={handleSubmit}>
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
                        disabled={loading}
                    />
                    <StyledTextField
                        margin="normal"
                        required
                        fullWidth
                        name="newPassword"
                        label="新密碼"
                        type={showPassword.newPassword ? "text" : "password"}
                        id="newPassword"
                        value={formData.newPassword}
                        onChange={handleChange}
                        disabled={loading}
                        InputProps={{
                            endAdornment: (
                                <IconButton
                                    onClick={() => setShowPassword({
                                        ...showPassword,
                                        newPassword: !showPassword.newPassword
                                    })}
                                    edge="end"
                                >
                                    {showPassword.newPassword ? <VisibilityOffIcon /> : <VisibilityIcon />}
                                </IconButton>
                            ),
                        }}
                    />
                    <StyledTextField
                        margin="normal"
                        required
                        fullWidth
                        name="confirmPassword"
                        label="確認新密碼"
                        type={showPassword.confirmPassword ? "text" : "password"}
                        id="confirmPassword"
                        value={formData.confirmPassword}
                        onChange={handleChange}
                        disabled={loading}
                        InputProps={{
                            endAdornment: (
                                <IconButton
                                    onClick={() => setShowPassword({
                                        ...showPassword,
                                        confirmPassword: !showPassword.confirmPassword
                                    })}
                                    edge="end"
                                >
                                    {showPassword.confirmPassword ? <VisibilityOffIcon /> : <VisibilityIcon />}
                                </IconButton>
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
                        {loading ? '處理中...' : '重設密碼'}
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
                            返回登入頁面
                        </Link>
                    </Box>
                </Box>
            </StyledPaper>
        </StyledContainer>
    );
};

export default ForgotPassword; 