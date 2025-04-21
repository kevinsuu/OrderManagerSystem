import React, { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
    AppBar,
    Box,
    Container,
    Toolbar,
    Typography,
    Button,
    IconButton,
    Badge,
    Menu,
    MenuItem,
    Avatar,
    Paper,
    InputBase,
} from '@mui/material';
import {
    ShoppingCart as ShoppingCartIcon,
    Person as PersonIcon,
    Search as SearchIcon,
} from '@mui/icons-material';

const StoreLayout = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [anchorEl, setAnchorEl] = useState(null);
    const [user, setUser] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        // 檢查是否已登入
        const checkLoginStatus = () => {
            const userToken = localStorage.getItem('userToken');
            const userData = localStorage.getItem('userData');

            if (userToken && userData) {
                setUser(JSON.parse(userData));
            } else {
                setUser(null);
            }
        };

        // 初始檢查
        checkLoginStatus();

        // 監聽 storage 變化
        const handleStorageChange = (e) => {
            if (e.key === 'userToken' || e.key === 'userData') {
                checkLoginStatus();
            }
        };

        window.addEventListener('storage', handleStorageChange);

        // 創建自定義事件來處理同一頁面的狀態更新
        const handleLoginStateChange = () => {
            checkLoginStatus();
        };

        window.addEventListener('loginStateChange', handleLoginStateChange);

        return () => {
            window.removeEventListener('storage', handleStorageChange);
            window.removeEventListener('loginStateChange', handleLoginStateChange);
        };
    }, []);

    // 從 URL 獲取搜尋詞
    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const query = params.get('query');
        const isHomePage = location.pathname === '/' || location.pathname === '/store';

        if (query) {
            setSearchTerm(decodeURIComponent(query));
        } else if (isHomePage) {
            setSearchTerm('');
        }
    }, [location]);

    const handleMenuOpen = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleMenuClose = () => {
        setAnchorEl(null);
    };

    const handleLogout = () => {
        localStorage.removeItem('userToken');
        localStorage.removeItem('userData');
        setUser(null);
        // 觸發登入狀態變更事件
        window.dispatchEvent(new Event('loginStateChange'));
        handleMenuClose();
        navigate('/');
    };

    // 統一的搜尋處理函數
    const executeSearch = (searchValue) => {
        if (searchValue.trim()) {
            navigate(`/store/products/search?query=${encodeURIComponent(searchValue.trim())}&page=1&limit=10`);
        }
    };

    const handleSearch = (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            executeSearch(searchTerm);
        }
    };

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            <AppBar position="fixed">
                <Toolbar sx={{ gap: 2 }}>
                    <Typography
                        variant="h6"
                        component="div"
                        sx={{ cursor: 'pointer' }}
                        onClick={() => navigate('/')}
                    >
                        訂單管理系統
                    </Typography>

                    {/* 搜尋欄 */}
                    <Paper
                        component="form"
                        onSubmit={(e) => e.preventDefault()}
                        sx={{
                            p: '2px 4px',
                            display: 'flex',
                            alignItems: 'center',
                            width: 400,
                            bgcolor: 'white',
                            borderRadius: 2,
                        }}
                    >
                        <InputBase
                            sx={{ ml: 1, flex: 1 }}
                            placeholder="搜尋商品..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            onKeyPress={handleSearch}
                        />
                        <IconButton
                            sx={{ p: '10px' }}
                            onClick={() => executeSearch(searchTerm)}
                        >
                            <SearchIcon />
                        </IconButton>
                    </Paper>

                    {/* 右側按鈕 */}
                    <Box sx={{ display: 'flex', gap: 2, alignItems: 'center', ml: 'auto' }}>
                        <IconButton color="inherit" onClick={() => navigate('/cart')}>
                            <Badge badgeContent={0} color="error">
                                <ShoppingCartIcon />
                            </Badge>
                        </IconButton>

                        {user ? (
                            <>
                                <IconButton
                                    onClick={handleMenuOpen}
                                    sx={{
                                        padding: 0.5,
                                        border: '2px solid white',
                                        borderRadius: '50%',
                                    }}
                                >
                                    <Avatar
                                        alt={user.name}
                                        src={user.avatar}
                                        sx={{ width: 32, height: 32 }}
                                    >
                                        {user.name?.charAt(0)}
                                    </Avatar>
                                </IconButton>
                                <Menu
                                    anchorEl={anchorEl}
                                    open={Boolean(anchorEl)}
                                    onClose={handleMenuClose}
                                >
                                    <MenuItem onClick={() => {
                                        handleMenuClose();
                                        navigate('/profile');
                                    }}>
                                        個人資料
                                    </MenuItem>
                                    <MenuItem onClick={() => {
                                        handleMenuClose();
                                        navigate('/orders');
                                    }}>
                                        訂單記錄
                                    </MenuItem>
                                    <MenuItem onClick={() => {
                                        handleMenuClose();
                                        navigate('/wishlist');
                                    }}>
                                        收藏清單
                                    </MenuItem>
                                    <MenuItem onClick={handleLogout}>
                                        登出
                                    </MenuItem>
                                </Menu>
                            </>
                        ) : (
                            <Button
                                color="inherit"
                                onClick={() => navigate('/login')}
                                startIcon={<PersonIcon />}
                            >
                                登入
                            </Button>
                        )}
                    </Box>
                </Toolbar>
            </AppBar>
            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    bgcolor: 'background.default',
                    mt: ['48px', '56px', '64px'],
                    minHeight: '100vh',
                }}
            >
                <Outlet />
            </Box>
            <Box
                component="footer"
                sx={{
                    py: 3,
                    px: 2,
                    mt: 'auto',
                    backgroundColor: (theme) =>
                        theme.palette.mode === 'light'
                            ? theme.palette.grey[200]
                            : theme.palette.grey[800],
                }}
            >
                <Container maxWidth="lg">
                    <Typography variant="body2" color="text.secondary" align="center">
                        © {new Date().getFullYear()} 訂單管理系統. All rights reserved.
                    </Typography>
                </Container>
            </Box>
        </Box>
    );
};

export default StoreLayout; 