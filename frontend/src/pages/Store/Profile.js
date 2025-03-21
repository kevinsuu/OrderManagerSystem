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
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'https://ordermanagersystem-auth-service.onrender.com';

const Profile = () => {
    const navigate = useNavigate();
    const [user, setUser] = useState(null);
    const [openEditDialog, setOpenEditDialog] = useState(false);
    const [editForm, setEditForm] = useState({
        username: '',
        email: '',
    });

    useEffect(() => {
        const userData = localStorage.getItem('userData');
        if (!userData) {
            navigate('/login');
            return;
        }
        setUser(JSON.parse(userData));
        setEditForm({
            username: JSON.parse(userData).username,
            email: JSON.parse(userData).email,
        });
    }, [navigate]);

    const handleLogout = () => {
        localStorage.removeItem('userToken');
        localStorage.removeItem('userData');
        navigate('/login');
    };

    const handleEditClick = () => {
        setOpenEditDialog(true);
    };

    const handleEditClose = () => {
        setOpenEditDialog(false);
    };

    const handleEditSubmit = async () => {
        try {
            const token = localStorage.getItem('userToken');
            const response = await axios.put(
                `${API_URL}/users/${user.id}`,
                editForm,
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );

            if (response.data) {
                const updatedUser = { ...user, ...editForm };
                localStorage.setItem('userData', JSON.stringify(updatedUser));
                setUser(updatedUser);
                setOpenEditDialog(false);
            }
        } catch (error) {
            console.error('更新個人資料失敗:', error);
        }
    };

    if (!user) {
        return null;
    }

    return (
        <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
            <Paper elevation={3} sx={{ p: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 4 }}>
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

                <List>
                    <ListItem>
                        <ListItemText
                            primary="會員角色"
                            secondary={user.role === 'admin' ? '管理員' : '一般會員'}
                        />
                    </ListItem>
                    <Divider />
                    <ListItem>
                        <ListItemText
                            primary="帳號狀態"
                            secondary={user.status === 'active' ? '正常' : '停用'}
                        />
                    </ListItem>
                    <Divider />
                    <ListItem>
                        <ListItemText
                            primary="註冊時間"
                            secondary={new Date(user.created_at).toLocaleString()}
                        />
                    </ListItem>
                </List>

                <Box sx={{ mt: 4, display: 'flex', gap: 2 }}>
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={handleEditClick}
                    >
                        編輯資料
                    </Button>
                    <Button
                        variant="outlined"
                        color="error"
                        onClick={handleLogout}
                    >
                        登出
                    </Button>
                </Box>
            </Paper>

            <Dialog open={openEditDialog} onClose={handleEditClose}>
                <DialogTitle>編輯個人資料</DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        margin="dense"
                        label="用戶名稱"
                        type="text"
                        fullWidth
                        value={editForm.username}
                        onChange={(e) =>
                            setEditForm({ ...editForm, username: e.target.value })
                        }
                    />
                    <TextField
                        margin="dense"
                        label="電子郵件"
                        type="email"
                        fullWidth
                        value={editForm.email}
                        onChange={(e) =>
                            setEditForm({ ...editForm, email: e.target.value })
                        }
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleEditClose}>取消</Button>
                    <Button onClick={handleEditSubmit} variant="contained">
                        儲存
                    </Button>
                </DialogActions>
            </Dialog>
        </Container>
    );
};

export default Profile; 