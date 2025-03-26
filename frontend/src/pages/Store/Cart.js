import React, { useState } from 'react';
import {
    Container,
    Paper,
    Typography,
    Box,
    IconButton,
    Button,
    Divider,
    Grid,
    Card,
    CardContent,
    TextField,
    Alert,
    Collapse,
    Fade,
} from '@mui/material';
import {
    Add as AddIcon,
    Remove as RemoveIcon,
    Delete as DeleteIcon,
    ShoppingCart as CartIcon,
    LocalShipping as ShippingIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

const Cart = () => {
    const navigate = useNavigate();
    const [cartItems, setCartItems] = useState([
        {
            id: 1,
            name: "商品範例 1",
            price: 1000,
            quantity: 1,
            image: "https://via.placeholder.com/150",
            stock: 10
        },
        {
            id: 2,
            name: "商品範例 2",
            price: 500,
            quantity: 2,
            image: "https://via.placeholder.com/150",
            stock: 5
        }
    ]);
    const [showAlert, setShowAlert] = useState(false);

    // 計算總金額
    const total = cartItems.reduce((sum, item) => sum + (item.price * item.quantity), 0);

    // 運費計算（訂單滿 2000 免運費）
    const shippingFee = total >= 2000 ? 0 : 100;

    // 更新商品數量
    const updateQuantity = (id, newQuantity) => {
        setCartItems(items =>
            items.map(item =>
                item.id === id
                    ? { ...item, quantity: Math.max(1, Math.min(newQuantity, item.stock)) }
                    : item
            )
        );
    };

    // 移除商品
    const removeItem = (id) => {
        setCartItems(items => items.filter(item => item.id !== id));
        setShowAlert(true);
        setTimeout(() => setShowAlert(false), 3000);
    };

    // 前往結帳
    const handleCheckout = () => {
        navigate('/checkout');
    };

    return (
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <CartIcon /> 購物車
            </Typography>

            <Grid container spacing={3}>
                <Grid item xs={12} md={8}>
                    <Collapse in={showAlert}>
                        <Alert
                            severity="success"
                            sx={{ mb: 2 }}
                            onClose={() => setShowAlert(false)}
                        >
                            商品已從購物車中移除
                        </Alert>
                    </Collapse>

                    {cartItems.length > 0 ? (
                        cartItems.map((item) => (
                            <Fade in={true} key={item.id}>
                                <Card sx={{ mb: 2, position: 'relative' }}>
                                    <CardContent>
                                        <Grid container spacing={2} alignItems="center">
                                            <Grid item xs={12} sm={3}>
                                                <img
                                                    src={item.image}
                                                    alt={item.name}
                                                    style={{
                                                        width: '100%',
                                                        height: 'auto',
                                                        borderRadius: '8px'
                                                    }}
                                                />
                                            </Grid>
                                            <Grid item xs={12} sm={9}>
                                                <Box sx={{
                                                    display: 'flex',
                                                    justifyContent: 'space-between',
                                                    alignItems: 'flex-start',
                                                    mb: 2
                                                }}>
                                                    <Typography variant="h6" component="div">
                                                        {item.name}
                                                    </Typography>
                                                    <IconButton
                                                        color="error"
                                                        onClick={() => removeItem(item.id)}
                                                        size="small"
                                                    >
                                                        <DeleteIcon />
                                                    </IconButton>
                                                </Box>
                                                <Box sx={{
                                                    display: 'flex',
                                                    justifyContent: 'space-between',
                                                    alignItems: 'center'
                                                }}>
                                                    <Box sx={{
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        gap: 1
                                                    }}>
                                                        <IconButton
                                                            size="small"
                                                            onClick={() => updateQuantity(item.id, item.quantity - 1)}
                                                            disabled={item.quantity <= 1}
                                                        >
                                                            <RemoveIcon />
                                                        </IconButton>
                                                        <TextField
                                                            size="small"
                                                            value={item.quantity}
                                                            onChange={(e) => {
                                                                const value = parseInt(e.target.value) || 1;
                                                                updateQuantity(item.id, value);
                                                            }}
                                                            inputProps={{
                                                                min: 1,
                                                                max: item.stock,
                                                                style: { textAlign: 'center', width: '50px' }
                                                            }}
                                                        />
                                                        <IconButton
                                                            size="small"
                                                            onClick={() => updateQuantity(item.id, item.quantity + 1)}
                                                            disabled={item.quantity >= item.stock}
                                                        >
                                                            <AddIcon />
                                                        </IconButton>
                                                    </Box>
                                                    <Typography variant="h6" color="primary">
                                                        NT$ {item.price * item.quantity}
                                                    </Typography>
                                                </Box>
                                                <Typography variant="caption" color="text.secondary">
                                                    庫存: {item.stock} 件
                                                </Typography>
                                            </Grid>
                                        </Grid>
                                    </CardContent>
                                </Card>
                            </Fade>
                        ))
                    ) : (
                        <Paper sx={{ p: 3, textAlign: 'center' }}>
                            <Typography variant="h6" color="text.secondary">
                                購物車是空的
                            </Typography>
                            <Button
                                variant="contained"
                                sx={{ mt: 2 }}
                                onClick={() => navigate('/products')}
                            >
                                繼續購物
                            </Button>
                        </Paper>
                    )}
                </Grid>

                <Grid item xs={12} md={4}>
                    <Paper sx={{ p: 3, position: 'sticky', top: 20 }}>
                        <Typography variant="h6" gutterBottom>
                            訂單摘要
                        </Typography>
                        <Box sx={{ my: 2 }}>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                                <Typography>商品總計</Typography>
                                <Typography>NT$ {total}</Typography>
                            </Box>
                            <Box sx={{
                                display: 'flex',
                                justifyContent: 'space-between',
                                mb: 1,
                                alignItems: 'center'
                            }}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <ShippingIcon fontSize="small" />
                                    <Typography>運費</Typography>
                                </Box>
                                <Typography>
                                    {shippingFee === 0 ? '免運費' : `NT$ ${shippingFee}`}
                                </Typography>
                            </Box>
                            {total < 2000 && (
                                <Typography variant="caption" color="primary">
                                    還差 NT$ {2000 - total} 即可享有免運優惠
                                </Typography>
                            )}
                        </Box>
                        <Divider sx={{ my: 2 }} />
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                            <Typography variant="h6">總計</Typography>
                            <Typography variant="h6" color="primary">
                                NT$ {total + shippingFee}
                            </Typography>
                        </Box>
                        <Button
                            variant="contained"
                            fullWidth
                            size="large"
                            onClick={handleCheckout}
                            disabled={cartItems.length === 0}
                        >
                            前往結帳
                        </Button>
                    </Paper>
                </Grid>
            </Grid>
        </Container>
    );
};

export default Cart; 