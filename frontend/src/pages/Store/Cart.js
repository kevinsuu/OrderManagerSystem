import React, { useState, useEffect, useCallback } from 'react';
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
import debounce from 'lodash/debounce';
import { createAuthAxios } from '../../utils/auth';

const CART_SERVICE_URL = 'https://ordermanagersystem.onrender.com';

const Cart = () => {
    const navigate = useNavigate();
    const [cartItems, setCartItems] = useState([]);
    const [showAlert, setShowAlert] = useState(false);
    const [alertMessage, setAlertMessage] = useState('');
    const [alertSeverity, setAlertSeverity] = useState('success');
    const [pendingUpdates, setPendingUpdates] = useState({});

    // 使用 useMemo 記憶化 authAxios 實例
    const authAxios = React.useMemo(() => createAuthAxios(navigate), [navigate]);

    // 成功訊息顯示
    const showSuccessMessage = useCallback((message) => {
        setShowAlert(true);
        setAlertMessage(message);
        setAlertSeverity('success');
        setTimeout(() => setShowAlert(false), 3000);
    }, []);

    // 錯誤訊息顯示
    const showErrorMessage = useCallback((message) => {
        setShowAlert(true);
        setAlertMessage(message);
        setAlertSeverity('error');
        setTimeout(() => setShowAlert(false), 3000);
    }, []);

    // 錯誤處理函數
    const handleError = useCallback((error) => {
        // 不需要在這裡處理 401，已在 axios 實例中處理
        showErrorMessage(error.response?.data?.error || '操作失敗');
    }, [showErrorMessage]);



    useEffect(() => {
        let isMounted = true;
        const controller = new AbortController();

        const getCartData = async () => {
            try {
                console.log('開始獲取購物車...');
                const response = await authAxios.get(`${CART_SERVICE_URL}/api/v1/cart/`, {
                    signal: controller.signal
                });

                if (!isMounted) return;

                const data = response.data;
                console.log('成功獲取購物車數據:', data);

                if (data.items) {
                    setCartItems(data.items.map(item => ({
                        id: item.ProductID,
                        name: item.Name,
                        price: item.Price,
                        quantity: item.Quantity,
                        image: item.Image || "https://via.placeholder.com/150",
                        stock: item.StockCount
                    })));
                }
            } catch (error) {
                if (error.name === 'AbortError' || !isMounted) return;

                console.error('獲取購物車失敗:', error);
                showErrorMessage(error.message);
            }
        };

        getCartData();

        return () => {
            isMounted = false;
            controller.abort();
        };
    }, [authAxios, showErrorMessage]);

    // 修正 debounce 使用方式
    const debouncedUpdateQuantity = useCallback((id, quantity) => {
        // 將 debounce 函數移到 useCallback 內部
        const updateQuantity = debounce(async (itemId, itemQuantity) => {
            try {
                await authAxios.put(`${CART_SERVICE_URL}/api/v1/cart/items`, {
                    ProductID: itemId,
                    Quantity: itemQuantity
                });

                // 清除待更新狀態
                setPendingUpdates(prev => {
                    const newUpdates = { ...prev };
                    delete newUpdates[itemId];
                    return newUpdates;
                });

                showSuccessMessage('數量已更新');
            } catch (error) {
                console.error('更新數量失敗:', error);
                handleError(error);
            }
        }, 1000);

        updateQuantity(id, quantity);
    }, [authAxios, showSuccessMessage, handleError]);

    // 處理數量變更
    const handleQuantityChange = (id, newQuantity) => {
        // 驗證輸入值
        const quantity = Math.max(1, Math.min(99, Number(newQuantity) || 1));

        // 更新本地狀態
        setCartItems(prev =>
            prev.map(item =>
                item.id === id ? { ...item, quantity } : item
            )
        );

        // 標記為待更新
        setPendingUpdates(prev => ({
            ...prev,
            [id]: quantity
        }));

        // 觸發防抖更新
        debouncedUpdateQuantity(id, quantity);
    };

    // 移除商品
    const removeItem = async (id) => {
        try {
            await authAxios.delete(`${CART_SERVICE_URL}/api/v1/cart/items/${id}`);
            setCartItems(items => items.filter(item => item.id !== id));
            showSuccessMessage('商品已從購物車中移除');
        } catch (error) {
            console.error('移除商品失敗:', error);
            handleError(error);
        }
    };

    // 清空購物車
    const handleClearCart = async () => {
        try {
            await authAxios.delete(`${CART_SERVICE_URL}/api/v1/cart/`);
            setCartItems([]);
            showSuccessMessage('購物車已清空');
        } catch (error) {
            console.error('清空購物車失敗:', error);
            handleError(error);
        }
    };

    // 計算總金額
    const total = cartItems.reduce((sum, item) => sum + (item.price * item.quantity), 0);

    // 運費計算（訂單滿 2000 免運費）
    const shippingFee = total >= 2000 ? 0 : 100;

    // 前往結帳
    const handleCheckout = () => {
        navigate('/checkout');
    };

    // 立即更新的函數
    const updateQuantityImmediately = useCallback(async (id, quantity) => {
        try {
            await authAxios.put(`${CART_SERVICE_URL}/api/v1/cart/items`, {
                ProductID: id,
                Quantity: quantity
            });
        } catch (error) {
            console.error('離開頁面時更新數量失敗:', error);
        }
    }, [authAxios]);

    // 在組件卸載前確保所有更新都被發送
    useEffect(() => {
        return () => {
            // 檢查所有待更新項目並立即發送
            Object.entries(pendingUpdates).forEach(([id, quantity]) => {
                updateQuantityImmediately(id, quantity);
            });
        };
    }, [pendingUpdates, updateQuantityImmediately]);

    return (
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Box sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                mb: 2
            }}>
                <Typography variant="h4" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <CartIcon /> 購物車
                </Typography>
                <Box sx={{ display: 'flex', gap: 2 }}>

                    {cartItems.length > 0 && (
                        <Button
                            variant="outlined"
                            color="error"
                            onClick={handleClearCart}
                            startIcon={<DeleteIcon />}
                        >
                            清空購物車
                        </Button>
                    )}
                </Box>
            </Box>

            <Grid container spacing={3}>
                <Grid item xs={12} md={8}>
                    <Collapse in={showAlert}>
                        <Alert
                            severity={alertSeverity}
                            sx={{ mb: 2 }}
                            onClose={() => setShowAlert(false)}
                        >
                            {alertMessage}
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
                                                    onError={(e) => {
                                                        console.error('圖片加載失敗:', item.image?.substring(0, 50));
                                                        e.target.src = "https://via.placeholder.com/150";
                                                    }}
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
                                                            onClick={() => handleQuantityChange(item.id, item.quantity - 1)}
                                                            disabled={item.quantity <= 1}
                                                        >
                                                            <RemoveIcon />
                                                        </IconButton>
                                                        <TextField
                                                            type="number"
                                                            value={item.quantity}
                                                            onChange={(e) => handleQuantityChange(item.id, e.target.value)}
                                                            inputProps={{
                                                                min: 1,
                                                                max: 99,
                                                                step: 1
                                                            }}
                                                            sx={{ width: '80px' }}
                                                            size="small"
                                                            variant="outlined"
                                                            onBlur={() => {
                                                                // 當輸入框失去焦點時立即更新
                                                                if (pendingUpdates[item.id]) {
                                                                    debouncedUpdateQuantity.flush(); // 如果使用lodash的debounce，可以立即執行
                                                                }
                                                            }}
                                                        />
                                                        <IconButton
                                                            size="small"
                                                            onClick={() => handleQuantityChange(item.id, item.quantity + 1)}
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
                                                {pendingUpdates[item.id] && (
                                                    <Typography variant="caption" color="text.secondary">
                                                        更新中...
                                                    </Typography>
                                                )}
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
                                onClick={() => navigate('/store/products')}
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