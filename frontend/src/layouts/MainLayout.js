import React from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import {
    AppBar,
    Box,
    CssBaseline,
    Drawer,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Toolbar,
    Typography,
    Button,
} from '@mui/material';
import {
    Dashboard as DashboardIcon,
    Inventory as InventoryIcon,
    Category as CategoryIcon,
    ShoppingCart as OrderIcon,
    LocalOffer as PromotionIcon,
    People as CustomerIcon,
    Assessment as ReportIcon,
    Logout as LogoutIcon,
} from '@mui/icons-material';

const drawerWidth = 280;

const menuItems = [
    { text: '儀表板', icon: <DashboardIcon />, path: '/admin/dashboard' },
    { text: '商品管理', icon: <InventoryIcon />, path: '/admin/products' },
    { text: '分類管理', icon: <CategoryIcon />, path: '/admin/categories' },
    { text: '訂單管理', icon: <OrderIcon />, path: '/admin/orders' },
    { text: '促銷活動', icon: <PromotionIcon />, path: '/admin/promotions' },
    { text: '會員管理', icon: <CustomerIcon />, path: '/admin/customers' },
    { text: '銷售報表', icon: <ReportIcon />, path: '/admin/reports' },
];

const MainLayout = () => {
    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem('adminLoggedIn');
        navigate('/admin/login');
    };

    return (
        <Box sx={{ display: 'flex' }}>
            <CssBaseline />
            <AppBar
                position="fixed"
                sx={{
                    width: `calc(100% - ${drawerWidth}px)`,
                    ml: `${drawerWidth}px`,
                }}
            >
                <Toolbar>
                    <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
                        訂單管理系統
                    </Typography>
                    <Button
                        color="inherit"
                        onClick={handleLogout}
                        startIcon={<LogoutIcon />}
                    >
                        登出
                    </Button>
                </Toolbar>
            </AppBar>
            <Drawer
                sx={{
                    width: drawerWidth,
                    flexShrink: 0,
                    '& .MuiDrawer-paper': {
                        width: drawerWidth,
                        boxSizing: 'border-box',
                        backgroundColor: 'background.paper',
                        borderRight: '1px solid rgba(0, 0, 0, 0.12)',
                    },
                }}
                variant="permanent"
                anchor="left"
            >
                <Toolbar
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        px: [1],
                        backgroundColor: 'primary.main',
                        color: 'white',
                    }}
                >
                    <Typography variant="h6" noWrap component="div">
                        管理後台
                    </Typography>
                </Toolbar>
                <List>
                    {menuItems.map((item) => (
                        <ListItem
                            button
                            key={item.text}
                            onClick={() => navigate(item.path)}
                            sx={{
                                '&:hover': {
                                    backgroundColor: 'action.hover',
                                },
                            }}
                        >
                            <ListItemIcon>{item.icon}</ListItemIcon>
                            <ListItemText primary={item.text} />
                        </ListItem>
                    ))}
                </List>
            </Drawer>
            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    bgcolor: 'background.default',
                    p: 3,
                    width: `calc(100% - ${drawerWidth}px)`,
                }}
            >
                <Toolbar />
                <Outlet />
            </Box>
        </Box>
    );
};

export default MainLayout; 