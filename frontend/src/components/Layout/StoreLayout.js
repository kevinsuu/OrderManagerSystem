import React, { useState } from 'react';
import {
    AppBar,
    Box,
    Toolbar,
    Typography,
    Button,
    IconButton,
    Badge,
    Menu,
    MenuItem,
    Container,
    InputBase,
    Drawer,
    List,
    ListItem,
    ListItemText,
    Divider,
} from '@mui/material';
import { styled, alpha } from '@mui/material/styles';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import PersonIcon from '@mui/icons-material/Person';
import SearchIcon from '@mui/icons-material/Search';
import MenuIcon from '@mui/icons-material/Menu';
import { useNavigate } from 'react-router-dom';

const Search = styled('div')(({ theme }) => ({
    position: 'relative',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: alpha(theme.palette.common.white, 0.15),
    '&:hover': {
        backgroundColor: alpha(theme.palette.common.white, 0.25),
    },
    marginRight: theme.spacing(2),
    marginLeft: 0,
    width: '100%',
    [theme.breakpoints.up('sm')]: {
        marginLeft: theme.spacing(3),
        width: 'auto',
    },
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
    padding: theme.spacing(0, 2),
    height: '100%',
    position: 'absolute',
    pointerEvents: 'none',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
    color: 'inherit',
    '& .MuiInputBase-input': {
        padding: theme.spacing(1, 1, 1, 0),
        paddingLeft: `calc(1em + ${theme.spacing(4)})`,
        transition: theme.transitions.create('width'),
        width: '100%',
        [theme.breakpoints.up('md')]: {
            width: '40ch',
        },
    },
}));

const categories = [
    '手機',
    '筆記型電腦',
    '平板電腦',
    '耳機',
    '智慧手錶',
    '配件',
];

const StoreLayout = ({ children }) => {
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
    const [anchorEl, setAnchorEl] = useState(null);
    const navigate = useNavigate();

    const handleProfileMenuOpen = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleMenuClose = () => {
        setAnchorEl(null);
    };

    const handleMobileMenuToggle = () => {
        setMobileMenuOpen(!mobileMenuOpen);
    };

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            <AppBar position="fixed" color="primary">
                <Container maxWidth="xl">
                    <Toolbar disableGutters>
                        <IconButton
                            color="inherit"
                            aria-label="open drawer"
                            edge="start"
                            onClick={handleMobileMenuToggle}
                            sx={{ mr: 2, display: { sm: 'none' } }}
                        >
                            <MenuIcon />
                        </IconButton>
                        <Typography
                            variant="h6"
                            noWrap
                            component="div"
                            sx={{ cursor: 'pointer' }}
                            onClick={() => navigate('/store')}
                        >
                            KS商城
                        </Typography>

                        <Box sx={{ display: { xs: 'none', sm: 'flex' }, ml: 4 }}>
                            {categories.map((category) => (
                                <Button
                                    key={category}
                                    color="inherit"
                                    onClick={() => navigate(`/store/category/${category}`)}
                                >
                                    {category}
                                </Button>
                            ))}
                        </Box>

                        <Box sx={{ flexGrow: 1 }} />

                        <Search>
                            <SearchIconWrapper>
                                <SearchIcon />
                            </SearchIconWrapper>
                            <StyledInputBase
                                placeholder="搜尋商品..."
                                inputProps={{ 'aria-label': 'search' }}
                            />
                        </Search>

                        <Box sx={{ display: 'flex', alignItems: 'center' }}>
                            <IconButton color="inherit" onClick={() => navigate('/store/cart')}>
                                <Badge badgeContent={3} color="error">
                                    <ShoppingCartIcon />
                                </Badge>
                            </IconButton>
                            <IconButton
                                edge="end"
                                color="inherit"
                                onClick={handleProfileMenuOpen}
                            >
                                <PersonIcon />
                            </IconButton>
                        </Box>
                    </Toolbar>
                </Container>
            </AppBar>

            <Drawer
                anchor="left"
                open={mobileMenuOpen}
                onClose={handleMobileMenuToggle}
            >
                <Box
                    sx={{ width: 250 }}
                    role="presentation"
                    onClick={handleMobileMenuToggle}
                >
                    <List>
                        {categories.map((category) => (
                            <ListItem
                                button
                                key={category}
                                onClick={() => navigate(`/store/category/${category}`)}
                            >
                                <ListItemText primary={category} />
                            </ListItem>
                        ))}
                    </List>
                </Box>
            </Drawer>

            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
                PaperProps={{
                    elevation: 0,
                    sx: {
                        filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.12))',
                        mt: 1.5,
                    },
                }}
            >
                <MenuItem onClick={() => { handleMenuClose(); navigate('/store/profile'); }}>
                    個人資料
                </MenuItem>
                <MenuItem onClick={() => { handleMenuClose(); navigate('/store/orders'); }}>
                    我的訂單
                </MenuItem>
                <Divider />
                <MenuItem onClick={handleMenuClose}>登出</MenuItem>
            </Menu>

            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    pt: { xs: 8, sm: 9 },
                    pb: 3,
                    backgroundColor: (theme) => theme.palette.grey[100],
                }}
            >
                {children}
            </Box>

            <Box
                component="footer"
                sx={{
                    py: 3,
                    px: 2,
                    mt: 'auto',
                    backgroundColor: (theme) => theme.palette.grey[200],
                }}
            >
                <Container maxWidth="lg">
                    <Typography variant="body2" color="text.secondary" align="center">
                        © {new Date().getFullYear()} KS商城. 版權所有.
                    </Typography>
                </Container>
            </Box>
        </Box>
    );
};

export default StoreLayout; 