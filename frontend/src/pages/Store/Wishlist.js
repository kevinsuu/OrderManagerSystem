import React, { useState, useEffect, useMemo, useCallback } from 'react';
import {
    Container,
    Typography,
    Box,
    Grid,
    Card,
    CardMedia,
    CardContent,
    CardActions,
    Button,
    IconButton,
    CircularProgress,
    Alert,
    Pagination,
    Snackbar,
} from '@mui/material';
import {
    Delete as DeleteIcon,
    ShoppingCart as CartIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { createAuthAxios } from '../../utils/auth';

const ITEMS_PER_PAGE = 12;

// 在生產環境中使用環境變數，在開發環境中使用備用值
const USER_API_URL = process.env.REACT_APP_USER_SERVICE_URL || 'http://localhost:8081';
const CART_API_URL = process.env.REACT_APP_CART_SERVICE_URL || 'http://localhost:8082';

const Wishlist = () => {
    const navigate = useNavigate();
    const [wishlist, setWishlist] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [snackbar, setSnackbar] = useState({
        open: false,
        message: '',
        severity: 'success'
    });

    const authAxios = useMemo(() => createAuthAxios(navigate), [navigate]);

    const fetchWishlist = useCallback(async () => {
        setLoading(true);
        try {
            const params = new URLSearchParams();
            params.set('page', page.toString());
            params.set('limit', ITEMS_PER_PAGE.toString());

            // 實際 API 呼叫
            const response = await authAxios.get(`${USER_API_URL}/api/v1/wishlist?${params.toString()}`);

            console.log('收藏清單API回應:', response.data);

            if (response.data.success) {
                setWishlist(response.data.data.wishlist || []);
                setTotal(response.data.data.total || 0);
                setError('');
            } else {
                setError(response.data.message || '獲取收藏清單失敗');
            }
        } catch (error) {
            console.error('獲取收藏清單失敗:', error);
            setError(error.response?.data?.message || '無法載入收藏清單');
        } finally {
            setLoading(false);
        }
    }, [authAxios, page]);

    useEffect(() => {
        fetchWishlist();
    }, [fetchWishlist]);

    const handlePageChange = (event, value) => {
        setPage(value);
        window.scrollTo(0, 0);
    };

    const handleRemoveFromWishlist = async (productId) => {
        try {
            // 實際 API 呼叫
            const response = await authAxios.delete(`${USER_API_URL}/api/v1/wishlist/${productId}`);

            if (response.data.success) {
                // 從列表中移除已刪除的項目
                setWishlist(wishlist.filter(item => (item.product.id || item.product._id) !== productId));
                setSnackbar({
                    open: true,
                    message: '已從收藏清單移除',
                    severity: 'success'
                });
            } else {
                setSnackbar({
                    open: true,
                    message: response.data.message || '移除失敗',
                    severity: 'error'
                });
            }
        } catch (error) {
            console.error('移除收藏失敗:', error);
            setSnackbar({
                open: true,
                message: error.response?.data?.message || '無法移除商品',
                severity: 'error'
            });
        }
    };

    const handleAddToCart = async (productId) => {
        try {
            // 實際 API 呼叫
            const response = await authAxios.post(`${CART_API_URL}/api/v1/cart/items`, {
                productId,
                quantity: 1
            });

            if (response.data.success) {
                setSnackbar({
                    open: true,
                    message: '已加入購物車',
                    severity: 'success'
                });
            } else {
                setSnackbar({
                    open: true,
                    message: response.data.message || '加入購物車失敗',
                    severity: 'error'
                });
            }
        } catch (error) {
            console.error('加入購物車失敗:', error);
            setSnackbar({
                open: true,
                message: error.response?.data?.message || '無法加入購物車',
                severity: 'error'
            });
        }
    };

    const handleCloseSnackbar = () => {
        setSnackbar({ ...snackbar, open: false });
    };

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Typography variant="h4" component="h1" gutterBottom>
                我的收藏
            </Typography>

            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            {loading ? (
                <Box display="flex" justifyContent="center" py={4}>
                    <CircularProgress />
                </Box>
            ) : wishlist.length === 0 ? (
                <Alert severity="info">
                    您目前沒有收藏任何商品
                </Alert>
            ) : (
                <Grid container spacing={2}>
                    {wishlist.map((item) => (
                        <Grid item key={item.id || item._id} xs={12} sm={6} md={4} lg={3}>
                            <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                                <CardMedia
                                    component="img"
                                    height="200"
                                    image={item.product.images?.[0]?.url || '/placeholder.png'}
                                    alt={item.product.name}
                                    sx={{ objectFit: 'contain', p: 2, cursor: 'pointer' }}
                                    onClick={() => navigate(`/store/products/${item.product.id || item.product._id}`)}
                                />
                                <CardContent sx={{ flexGrow: 1 }}>
                                    <Typography gutterBottom variant="h6" component="h2">
                                        {item.product.name}
                                    </Typography>
                                    <Typography variant="h6" color="primary" gutterBottom>
                                        NT$ {item.product.price.toLocaleString()}
                                    </Typography>
                                </CardContent>
                                <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
                                    <IconButton
                                        color="error"
                                        onClick={() => handleRemoveFromWishlist(item.product.id || item.product._id)}
                                    >
                                        <DeleteIcon />
                                    </IconButton>
                                    <Button
                                        variant="contained"
                                        startIcon={<CartIcon />}
                                        onClick={() => handleAddToCart(item.product.id || item.product._id)}
                                    >
                                        加入購物車
                                    </Button>
                                </CardActions>
                            </Card>
                        </Grid>
                    ))}
                </Grid>
            )}

            {total > ITEMS_PER_PAGE && (
                <Box sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
                    <Pagination
                        count={Math.ceil(total / ITEMS_PER_PAGE)}
                        page={page}
                        onChange={handlePageChange}
                        color="primary"
                    />
                </Box>
            )}

            <Snackbar
                open={snackbar.open}
                autoHideDuration={3000}
                onClose={handleCloseSnackbar}
                message={snackbar.message}
            />
        </Container>
    );
};

export default Wishlist; 