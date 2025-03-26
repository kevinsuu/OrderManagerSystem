import React, { useState, useEffect } from 'react';
import {
    Container,
    Paper,
    Typography,
    Box,
    Avatar,
    List,
    ListItem,
    ListItemText,
    Divider,
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    Grid,
    Switch,
    FormControlLabel,
    IconButton,
    Card,
    CardContent,
    ListItemIcon,
    Chip,
    Snackbar,
    Alert,
    Menu,
    MenuItem,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import SaveIcon from '@mui/icons-material/Save';
import AddIcon from '@mui/icons-material/Add';
import LogoutIcon from '@mui/icons-material/Logout';
import SettingsIcon from '@mui/icons-material/Settings';
import PersonIcon from '@mui/icons-material/Person';
import VerifiedUserIcon from '@mui/icons-material/VerifiedUser';
import CalendarTodayIcon from '@mui/icons-material/CalendarToday';
import LanguageIcon from '@mui/icons-material/Language';
import LocationOnIcon from '@mui/icons-material/LocationOn';
import LocationOffIcon from '@mui/icons-material/LocationOff';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'https://ordermanagersystem-auth-service.onrender.com';

const Profile = () => {
    const navigate = useNavigate();
    const [user, setUser] = useState(null);
    const [addresses, setAddresses] = useState([]);
    const [preferencesForm, setPreferencesForm] = useState({
        language: 'zh_TW',
        currency: 'TWD',
        theme: 'light',
        notification_email: true,
        notification_sms: false
    });
    const [openAddressDialog, setOpenAddressDialog] = useState(false);
    const [addressForm, setAddressForm] = useState({
        name: '',
        phone: '',
        postal_code: '',
        city: '',
        district: '',
        street: '',
        is_default: false
    });
    const [editingAddressIndex, setEditingAddressIndex] = useState(null);
    const [snackbar, setSnackbar] = useState({
        open: false,
        message: '',
        severity: 'success'
    });
    const [anchorEl, setAnchorEl] = useState(null);
    const [selectedAddress, setSelectedAddress] = useState(null);
    const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);

    useEffect(() => {
        const userData = localStorage.getItem('userData');
        if (!userData) {
            navigate('/login');
            return;
        }
        setUser(JSON.parse(userData));
        fetchUserData();
    }, [navigate]);

    const fetchUserData = async () => {
        try {
            const token = localStorage.getItem('userToken');
            const headers = {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            };

            // 獲取用戶偏好設定
            const preferencesResponse = await axios.get(
                `${AUTH_SERVICE_URL}/api/v1/user/preferences`,
                {
                    headers,
                    withCredentials: true  // 如果需要發送 cookies
                }
            );
            setPreferencesForm(preferencesResponse.data);

            // 獲取用戶地址
            const addressesResponse = await axios.get(
                `${AUTH_SERVICE_URL}/api/v1/user/addresses`,
                {
                    headers,
                    withCredentials: true  // 如果需要發送 cookies
                }
            );
            setAddresses(addressesResponse.data);
        } catch (error) {
            console.error('獲取用戶資料失敗:', error);
            // 添加更詳細的錯誤處理
            if (error.response) {
                console.error('錯誤狀態:', error.response.status);
                console.error('錯誤數據:', error.response.data);
            }
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('userToken');
        localStorage.removeItem('userData');
        navigate('/login');
    };

    const handlePreferencesSubmit = async () => {
        try {
            const token = localStorage.getItem('userToken');
            const response = await axios.put(
                `${AUTH_SERVICE_URL}/api/v1/user/preferences`,
                preferencesForm,
                {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                }
            );

            if (response.data) {
                setSnackbar({
                    open: true,
                    message: '偏好設定已更新成功！',
                    severity: 'success'
                });
            }
        } catch (error) {
            console.error('更新偏好設定失敗:', error);
            setSnackbar({
                open: true,
                message: '更新偏好設定失敗，請稍後再試',
                severity: 'error'
            });
        }
    };

    const handleAddressClick = (address, index) => {
        setAddressForm(address);
        setEditingAddressIndex(index);
        setOpenAddressDialog(true);
    };

    const handleAddNewAddress = () => {
        setAddressForm({
            name: '',
            phone: '',
            postal_code: '',
            city: '',
            district: '',
            street: '',
            is_default: false
        });
        setEditingAddressIndex(null);
        setOpenAddressDialog(true);
    };

    const handleAddressClose = () => {
        setOpenAddressDialog(false);
        setEditingAddressIndex(null);
    };

    const handleAddressMenuClick = (event, address) => {
        event.stopPropagation();
        setAnchorEl(event.currentTarget);
        setSelectedAddress(address);
    };

    const handleAddressMenuClose = () => {
        setAnchorEl(null);
        setSelectedAddress(null);
    };

    const handleEditAddress = () => {
        handleAddressMenuClose();
        const addressIndex = addresses.findIndex(addr => addr.id === selectedAddress.id);
        handleAddressClick(selectedAddress, addressIndex);
    };

    const handleDeleteAddress = async () => {
        try {
            const token = localStorage.getItem('userToken');
            await axios.delete(
                `${AUTH_SERVICE_URL}/api/v1/user/addresses/${selectedAddress.id}`,
                {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }
                }
            );

            await fetchUserData();
            setSnackbar({
                open: true,
                message: '地址已成功刪除！',
                severity: 'success'
            });
        } catch (error) {
            console.error('刪除地址失敗:', error);
            setSnackbar({
                open: true,
                message: '刪除地址失敗，請稍後再試',
                severity: 'error'
            });
        }
        handleAddressMenuClose();
        setDeleteConfirmOpen(false);
    };

    const handleAddressSubmit = async () => {
        try {
            const token = localStorage.getItem('userToken');
            const headers = {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            };

            let response;
            if (editingAddressIndex !== null) {
                // 更新現有地址
                response = await axios.put(
                    `${AUTH_SERVICE_URL}/api/v1/user/addresses/${addresses[editingAddressIndex].id}`,
                    addressForm,
                    { headers }
                );
            } else {
                // 新增地址
                response = await axios.post(
                    `${AUTH_SERVICE_URL}/api/v1/user/addresses`,
                    addressForm,
                    { headers }
                );
            }

            if (response.data) {
                await fetchUserData();
                setOpenAddressDialog(false);
                setSnackbar({
                    open: true,
                    message: `地址已成功${editingAddressIndex !== null ? '更新' : '新增'}！`,
                    severity: 'success'
                });
            }
        } catch (error) {
            console.error('更新地址失敗:', error);
            setSnackbar({
                open: true,
                message: `${editingAddressIndex !== null ? '更新' : '新增'}地址失敗，請稍後再試`,
                severity: 'error'
            });
        }
    };

    if (!user) {
        return null;
    }

    return (
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
                <Box sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'flex-start',
                    mb: 4
                }}>
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Avatar
                            sx={{
                                width: 100,
                                height: 100,
                                mr: 3,
                                bgcolor: 'primary.main',
                                fontSize: '2rem',
                            }}
                        >
                            {user.username?.[0]?.toUpperCase()}
                        </Avatar>
                        <Box>
                            <Typography variant="h4" gutterBottom>
                                {user.username}
                            </Typography>
                            <Typography variant="subtitle1" color="text.secondary">
                                {user.email}
                            </Typography>
                        </Box>
                    </Box>

                    <IconButton
                        color="error"
                        onClick={handleLogout}
                        sx={{
                            bgcolor: 'action.hover',
                            '&:hover': { bgcolor: 'error.light' }
                        }}
                    >
                        <LogoutIcon />
                    </IconButton>
                </Box>

                <List sx={{ bgcolor: 'background.paper', borderRadius: 1 }}>
                    <ListItem>
                        <ListItemIcon>
                            <PersonIcon />
                        </ListItemIcon>
                        <ListItemText
                            primary="會員角色"
                            secondary={user.role === 'admin' ? '管理員' : '一般會員'}
                        />
                    </ListItem>
                    <Divider component="li" />
                    <ListItem>
                        <ListItemIcon>
                            <VerifiedUserIcon />
                        </ListItemIcon>
                        <ListItemText
                            primary="帳號狀態"
                            secondary={user.status === 'active' ? '正常' : '停用'}
                        />
                    </ListItem>
                    <Divider component="li" />
                    <ListItem>
                        <ListItemIcon>
                            <CalendarTodayIcon />
                        </ListItemIcon>
                        <ListItemText
                            primary="註冊時間"
                            secondary={new Date(user.created_at).toLocaleString()}
                        />
                    </ListItem>
                </List>
            </Paper>

            <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
                <Box sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    mb: 2
                }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <SettingsIcon color="primary" />
                        <Typography variant="h5" sx={{ fontWeight: 'medium' }}>
                            個人化設定
                        </Typography>
                    </Box>
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={handlePreferencesSubmit}
                        startIcon={<SaveIcon />}
                        size="small"
                    >
                        儲存設定
                    </Button>
                </Box>

                <Grid container spacing={2}>
                    {[
                        {
                            icon: <LanguageIcon />,
                            title: '語言',
                            value: preferencesForm.language,
                            options: [
                                { value: 'zh_TW', label: '繁體中文' },
                                { value: 'en', label: 'English' }
                            ],
                            onChange: (e) => setPreferencesForm({ ...preferencesForm, language: e.target.value })
                        },
                        {
                            icon: <LanguageIcon />,
                            title: '幣別',
                            value: preferencesForm.currency,
                            options: [
                                { value: 'TWD', label: '新台幣 (TWD)' },
                                { value: 'USD', label: '美元 (USD)' }
                            ],
                            onChange: (e) => setPreferencesForm({ ...preferencesForm, currency: e.target.value })
                        },
                        {
                            icon: <LanguageIcon />,
                            title: '主題',
                            value: preferencesForm.theme,
                            options: [
                                { value: 'light', label: '淺色' },
                                { value: 'dark', label: '深色' }
                            ],
                            onChange: (e) => setPreferencesForm({ ...preferencesForm, theme: e.target.value })
                        },
                        {
                            icon: <LanguageIcon />,
                            title: '通知設定',
                            value: preferencesForm.notification_email.toString(),
                            options: [
                                { value: 'true', label: '啟用' },
                                { value: 'false', label: '停用' }
                            ],
                            onChange: (e) => {
                                const newValue = e.target.value === 'true';
                                setPreferencesForm(prev => ({
                                    ...prev,
                                    notification_email: newValue
                                }));
                            }
                        },
                    ].map((setting, index) => (
                        <Grid item xs={12} sm={6} md={3} key={index}>
                            <Card sx={{ height: '100%' }}>
                                <CardContent>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                                        {setting.icon}
                                        <Typography variant="subtitle1">
                                            {setting.title}
                                        </Typography>
                                    </Box>
                                    <TextField
                                        select
                                        fullWidth
                                        size="small"
                                        value={setting.value}
                                        onChange={setting.onChange}
                                        SelectProps={{ native: true }}
                                    >
                                        {setting.options.map(option => (
                                            <option key={option.value} value={option.value}>
                                                {option.label}
                                            </option>
                                        ))}
                                    </TextField>
                                </CardContent>
                            </Card>
                        </Grid>
                    ))}
                </Grid>
            </Paper>

            <Paper elevation={3} sx={{ p: 3 }}>
                <Box sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    mb: 2
                }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <LocationOnIcon color="primary" />
                        <Typography variant="h5" sx={{ fontWeight: 'medium' }}>
                            收貨地址
                        </Typography>
                    </Box>
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={handleAddNewAddress}
                        startIcon={<AddIcon />}
                        size="small"
                    >
                        新增地址
                    </Button>
                </Box>

                <Grid container spacing={2}>
                    {addresses.length > 0 ? (
                        addresses.map((address, index) => (
                            <Grid item xs={12} sm={6} md={4} key={index}>
                                <Card
                                    sx={{
                                        cursor: 'pointer',
                                        transition: 'all 0.3s',
                                        '&:hover': {
                                            transform: 'translateY(-4px)',
                                            boxShadow: 3
                                        }
                                    }}
                                >
                                    <CardContent>
                                        <Box sx={{
                                            display: 'flex',
                                            justifyContent: 'space-between',
                                            alignItems: 'center',
                                            mb: 1
                                        }}>
                                            <Typography variant="h6">
                                                {address.name}
                                            </Typography>
                                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                {address.is_default && (
                                                    <Chip
                                                        label="預設"
                                                        color="primary"
                                                        size="small"
                                                    />
                                                )}
                                                <IconButton
                                                    size="small"
                                                    onClick={(e) => handleAddressMenuClick(e, address)}
                                                >
                                                    <MoreVertIcon />
                                                </IconButton>
                                            </Box>
                                        </Box>
                                        <Typography variant="body2" sx={{ mb: 1 }}>
                                            {address.phone}
                                        </Typography>
                                        <Typography color="text.secondary" variant="body2">
                                            {address.postal_code} {address.city}{address.district}{address.street}
                                        </Typography>
                                    </CardContent>
                                </Card>
                            </Grid>
                        ))
                    ) : (
                        <Grid item xs={12}>
                            <Box
                                sx={{
                                    textAlign: 'center',
                                    py: 4,
                                    bgcolor: 'background.paper',
                                    borderRadius: 1
                                }}
                            >
                                <LocationOffIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
                                <Typography color="text.secondary">
                                    尚未新增地址
                                </Typography>
                            </Box>
                        </Grid>
                    )}
                </Grid>
            </Paper>

            <Dialog open={openAddressDialog} onClose={handleAddressClose} maxWidth="sm" fullWidth>
                <DialogTitle>
                    {editingAddressIndex !== null ? '編輯地址' : '新增地址'}
                </DialogTitle>
                <DialogContent>
                    <Box sx={{ mt: 2 }}>
                        <Grid container spacing={2}>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    required
                                    margin="dense"
                                    label="地址名稱"
                                    placeholder="例：公司、家"
                                    fullWidth
                                    value={addressForm.name}
                                    onChange={(e) =>
                                        setAddressForm({ ...addressForm, name: e.target.value })
                                    }
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    required
                                    margin="dense"
                                    label="聯絡電話"
                                    fullWidth
                                    value={addressForm.phone}
                                    onChange={(e) =>
                                        setAddressForm({ ...addressForm, phone: e.target.value })
                                    }
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    required
                                    margin="dense"
                                    label="郵遞區號"
                                    fullWidth
                                    value={addressForm.postal_code}
                                    onChange={(e) =>
                                        setAddressForm({ ...addressForm, postal_code: e.target.value })
                                    }
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    required
                                    margin="dense"
                                    label="城市"
                                    fullWidth
                                    value={addressForm.city}
                                    onChange={(e) =>
                                        setAddressForm({ ...addressForm, city: e.target.value })
                                    }
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    required
                                    margin="dense"
                                    label="區域"
                                    fullWidth
                                    value={addressForm.district}
                                    onChange={(e) =>
                                        setAddressForm({ ...addressForm, district: e.target.value })
                                    }
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    required
                                    margin="dense"
                                    label="街道地址"
                                    fullWidth
                                    value={addressForm.street}
                                    onChange={(e) =>
                                        setAddressForm({ ...addressForm, street: e.target.value })
                                    }
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={addressForm.is_default}
                                            onChange={(e) =>
                                                setAddressForm({
                                                    ...addressForm,
                                                    is_default: e.target.checked,
                                                })
                                            }
                                        />
                                    }
                                    label="設為預設地址"
                                />
                            </Grid>
                        </Grid>
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleAddressClose}>取消</Button>
                    <Button onClick={handleAddressSubmit} variant="contained">
                        儲存
                    </Button>
                </DialogActions>
            </Dialog>

            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleAddressMenuClose}
            >
                <MenuItem onClick={handleEditAddress}>
                    <EditIcon fontSize="small" sx={{ mr: 1 }} />
                    編輯地址
                </MenuItem>
                <MenuItem
                    onClick={() => setDeleteConfirmOpen(true)}
                    sx={{ color: 'error.main' }}
                >
                    <DeleteIcon fontSize="small" sx={{ mr: 1 }} />
                    刪除地址
                </MenuItem>
            </Menu>

            <Dialog
                open={deleteConfirmOpen}
                onClose={() => setDeleteConfirmOpen(false)}
            >
                <DialogTitle>確認刪除地址</DialogTitle>
                <DialogContent>
                    <Typography>
                        確定要刪除這個地址嗎？此操作無法復原。
                    </Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteConfirmOpen(false)}>
                        取消
                    </Button>
                    <Button
                        onClick={handleDeleteAddress}
                        color="error"
                        variant="contained"
                    >
                        刪除
                    </Button>
                </DialogActions>
            </Dialog>

            <Snackbar
                open={snackbar.open}
                autoHideDuration={2000}
                onClose={() => setSnackbar({ ...snackbar, open: false })}
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
            >
                <Alert
                    onClose={() => setSnackbar({ ...snackbar, open: false })}
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </Container>
    );
};

export default Profile; 