import React from 'react';
import {
    Box,
    AppBar,
    Toolbar,
    Typography,
    Drawer,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    IconButton,
    Badge,
    Avatar,
    Menu,
    MenuItem,
    Divider
} from '@mui/material';
import { styled } from '@mui/material/styles';
import MenuIcon from '@mui/icons-material/Menu';
import DashboardIcon from '@mui/icons-material/Dashboard';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import InventoryIcon from '@mui/icons-material/Inventory';
import PeopleIcon from '@mui/icons-material/People';
import CategoryIcon from '@mui/icons-material/Category';
import LocalOfferIcon from '@mui/icons-material/LocalOffer';
import ReceiptIcon from '@mui/icons-material/Receipt';
import AssessmentIcon from '@mui/icons-material/Assessment';
import NotificationsIcon from '@mui/icons-material/Notifications';
import { useNavigate } from 'react-router-dom';

const drawerWidth = 260;

const Main = styled('main', { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
        flexGrow: 1,
        padding: theme.spacing(3),
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        marginLeft: `-${drawerWidth}px`,
        ...(open && {
            transition: theme.transitions.create('margin', {
                easing: theme.transitions.easing.easeOut,
                duration: theme.transitions.duration.enteringScreen,
            }),
            marginLeft: 0,
        }),
        backgroundColor: theme.palette.background.default,
    }),
);

const StyledAppBar = styled(AppBar)(({ theme }) => ({
    backgroundColor: 'white',
    color: theme.palette.text.primary,
    boxShadow: '0 1px 3px rgba(0,0,0,0.12)',
}));

const StyledToolbar = styled(Toolbar)(({ theme }) => ({
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: theme.spacing(0, 2),
}));

const LogoSection = styled(Box)(({ theme }) => ({
    display: 'flex',
    alignItems: 'center',
    gap: theme.spacing(2),
}));

const ActionSection = styled(Box)(({ theme }) => ({
    display: 'flex',
    alignItems: 'center',
    gap: theme.spacing(2),
}));

const StyledBadge = styled(Badge)(({ theme }) => ({
    '& .MuiBadge-badge': {
        backgroundColor: theme.palette.error.main,
        color: 'white',
    },
}));

const MainLayout = ({ children }) => {
    const [open, setOpen] = React.useState(true);
    const [anchorEl, setAnchorEl] = React.useState(null);
    const navigate = useNavigate();

    const menuItems = [
        { text: '儀表板', icon: <DashboardIcon />, path: '/' },
        { text: '商品管理', icon: <InventoryIcon />, path: '/products' },
        { text: '商品分類', icon: <CategoryIcon />, path: '/categories' },
        { text: '訂單管理', icon: <ReceiptIcon />, path: '/orders' },
        { text: '促銷活動', icon: <LocalOfferIcon />, path: '/promotions' },
        { text: '會員管理', icon: <PeopleIcon />, path: '/customers' },
        { text: '銷售報表', icon: <AssessmentIcon />, path: '/reports' },
    ];

    const handleDrawerToggle = () => {
        setOpen(!open);
    };

    const handleProfileMenuOpen = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleProfileMenuClose = () => {
        setAnchorEl(null);
    };

    return (
        <Box sx={{ display: 'flex' }}>
            <StyledAppBar position="fixed">
                <StyledToolbar>
                    <LogoSection>
                        <IconButton
                            color="inherit"
                            aria-label="open drawer"
                            onClick={handleDrawerToggle}
                            edge="start"
                        >
                            <MenuIcon />
                        </IconButton>
                        <Typography variant="h6" noWrap component="div" sx={{ fontWeight: 600 }}>
                            商城管理系統
                        </Typography>
                    </LogoSection>
                    <ActionSection>
                        <IconButton color="inherit">
                            <StyledBadge badgeContent={4}>
                                <NotificationsIcon />
                            </StyledBadge>
                        </IconButton>
                        <IconButton color="inherit">
                            <StyledBadge badgeContent={2}>
                                <ShoppingCartIcon />
                            </StyledBadge>
                        </IconButton>
                        <IconButton onClick={handleProfileMenuOpen}>
                            <Avatar sx={{ width: 32, height: 32 }} />
                        </IconButton>
                        <Menu
                            anchorEl={anchorEl}
                            open={Boolean(anchorEl)}
                            onClose={handleProfileMenuClose}
                            PaperProps={{
                                elevation: 0,
                                sx: {
                                    filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.12))',
                                    mt: 1.5,
                                },
                            }}
                        >
                            <MenuItem onClick={handleProfileMenuClose}>個人資料</MenuItem>
                            <MenuItem onClick={handleProfileMenuClose}>帳號設定</MenuItem>
                            <Divider />
                            <MenuItem onClick={handleProfileMenuClose}>登出</MenuItem>
                        </Menu>
                    </ActionSection>
                </StyledToolbar>
            </StyledAppBar>
            <Drawer
                variant="persistent"
                anchor="left"
                open={open}
                sx={{
                    width: drawerWidth,
                    flexShrink: 0,
                    '& .MuiDrawer-paper': {
                        width: drawerWidth,
                        boxSizing: 'border-box',
                        borderRight: '1px solid rgba(0, 0, 0, 0.08)',
                        backgroundColor: '#ffffff',
                    },
                }}
            >
                <Toolbar />
                <Box sx={{ overflow: 'auto', mt: 2 }}>
                    <List>
                        {menuItems.map((item) => (
                            <ListItem
                                button
                                key={item.text}
                                onClick={() => navigate(item.path)}
                                sx={{
                                    mx: 2,
                                    borderRadius: '8px',
                                    mb: 1,
                                    '&:hover': {
                                        backgroundColor: 'rgba(25, 118, 210, 0.08)',
                                    },
                                }}
                            >
                                <ListItemIcon sx={{ minWidth: 40 }}>{item.icon}</ListItemIcon>
                                <ListItemText primary={item.text} />
                            </ListItem>
                        ))}
                    </List>
                </Box>
            </Drawer>
            <Main open={open}>
                <Toolbar />
                {children}
            </Main>
        </Box>
    );
};

export default MainLayout; 