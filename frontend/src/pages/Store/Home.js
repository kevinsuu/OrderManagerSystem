import React, { useState, useEffect, useCallback } from 'react';
import {
    Container,
    Grid,
    Card,
    CardMedia,
    CardContent,
    Typography,
    Box,
    Button,
    Rating,
    Paper,
    IconButton,
    Snackbar,
    Alert,
    Skeleton,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import { useNavigate } from 'react-router-dom';
import { createAuthAxios } from '../../utils/auth';

const StyledCard = styled(Card)(({ theme }) => ({
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    transition: 'all 0.3s ease-in-out',
    cursor: 'pointer',
    '&:hover': {
        transform: 'translateY(-4px)',
        boxShadow: '0 12px 24px -10px rgba(0,0,0,0.1)',
    },
}));

const ProductImage = styled(CardMedia)({
    paddingTop: '100%', // 1:1 比例
    position: 'relative',
    backgroundColor: '#f8fafc',
});

const DiscountBadge = styled(Box)(({ theme }) => ({
    position: 'absolute',
    top: 16,
    right: 16,
    backgroundColor: theme.palette.error.main,
    color: 'white',
    padding: '4px 8px',
    borderRadius: '4px',
    fontWeight: 600,
    fontSize: '0.875rem',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
}));

const ProductActions = styled(Box)(({ theme }) => ({
    position: 'absolute',
    bottom: 16,
    right: 16,
    display: 'flex',
    gap: theme.spacing(1),
    opacity: 0,
    transition: 'opacity 0.3s ease-in-out',
    '& .MuiIconButton-root': {
        backgroundColor: 'white',
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
        '&:hover': {
            backgroundColor: theme.palette.grey[100],
        },
    },
    '.MuiCard-root:hover &': {
        opacity: 1,
    },
}));



// 添加環境變數
const PRODUCT_SERVICE_URL = process.env.REACT_APP_PRODUCT_SERVICE_URL || 'https://ordermanagersystem-product-service.onrender.com';
const CART_SERVICE_URL = process.env.REACT_APP_CART_SERVICE_URL || 'https://ordermanagersystem.onrender.com';
const USER_SERVICE_URL = process.env.REACT_APP_USER_SERVICE_URL || 'https://ordermanagersystem.onrender.com';

// 添加骨架屏組件
const ProductSkeleton = () => (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
        <Skeleton variant="rectangular" sx={{ paddingTop: '100%' }} />
        <CardContent sx={{ flexGrow: 1 }}>
            <Skeleton variant="text" height={32} width="80%" sx={{ mb: 1 }} />
            <Skeleton variant="text" height={20} sx={{ mb: 1 }} />
            <Skeleton variant="text" height={20} width="60%" sx={{ mb: 2 }} />
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Skeleton variant="text" width={120} height={24} />
            </Box>
            <Skeleton variant="text" height={28} width="40%" />
        </CardContent>
    </Card>
);

// 優化產品卡片組件
const ProductCard = React.memo(({ product, isFavorite, onToggleFavorite, onAddToCart, onNavigate }) => {
    return (
        <StyledCard onClick={() => onNavigate(`/store/products/${product.id}`)}>
            <ProductImage
                image={product.images && product.images[0]?.url ? product.images[0].url : 'default-image-url.jpg'}
                title={product.name}
                loading="lazy"
            >
                {product.discount && (
                    <DiscountBadge>
                        {product.discount * 100}% OFF
                    </DiscountBadge>
                )}
                <ProductActions>
                    <IconButton
                        size="small"
                        color={isFavorite ? "error" : "default"}
                        onClick={(e) => onToggleFavorite(e, product.id)}
                    >
                        {isFavorite ?
                            <FavoriteIcon fontSize="small" /> :
                            <FavoriteBorderIcon fontSize="small" />
                        }
                    </IconButton>
                    <IconButton
                        size="small"
                        onClick={(e) => onAddToCart(e, product.id)}
                    >
                        <ShoppingCartIcon fontSize="small" />
                    </IconButton>
                </ProductActions>
            </ProductImage>
            <CardContent sx={{ flexGrow: 1 }}>
                <Typography
                    gutterBottom
                    variant="h6"
                    component="h2"
                    sx={{ fontWeight: 600 }}
                >
                    {product.name}
                </Typography>
                <Typography
                    variant="body2"
                    color="text.secondary"
                    sx={{ mb: 2 }}
                >
                    {product.description}
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <Rating value={product.rating || 0} precision={0.1} readOnly size="small" />
                    <Typography variant="body2" color="text.secondary" sx={{ ml: 1 }}>
                        ({product.reviews || 0})
                    </Typography>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'baseline', gap: 1 }}>
                    <Typography variant="h6" color="primary" sx={{ fontWeight: 600 }}>
                        ${product.price?.toLocaleString()}
                    </Typography>
                </Box>
            </CardContent>
        </StyledCard>
    );
});

const StorePage = () => {
    const navigate = useNavigate();
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [currentProductIndex, setCurrentProductIndex] = useState(0);
    const [featuredProducts, setFeaturedProducts] = useState([]);
    const [favorites, setFavorites] = useState({});
    const [alert, setAlert] = useState({ open: false, message: '', severity: 'success' });

    // 創建授權axios客戶端
    const authAxios = React.useMemo(() => createAuthAxios(navigate), [navigate]);

    // 顯示提示消息
    const showAlert = useCallback((message, severity = 'success') => {
        setAlert({
            open: true,
            message,
            severity
        });
    }, []);

    // 關閉提示
    const handleCloseAlert = useCallback(() => {
        setAlert(prev => ({ ...prev, open: false }));
    }, []);

    // 獲取收藏清單 - 非同步加載，不阻塞頁面顯示
    const fetchWishlist = useCallback(async () => {
        try {
            const response = await authAxios.get(`${USER_SERVICE_URL}/api/v1/wishlist/`);
            if (response.data && response.data.success && response.data.data && response.data.data.wishlist) {
                const wishlistItems = response.data.data.wishlist;
                const favoritesMap = {};
                wishlistItems.forEach(item => {
                    favoritesMap[item.productId] = true;
                });
                setFavorites(favoritesMap);
            }
        } catch (error) {
            console.error('Error fetching wishlist:', error);
        }
    }, [authAxios]);

    // 先加載商品數據
    useEffect(() => {
        const fetchProducts = async () => {
            try {
                const response = await authAxios.get(`${PRODUCT_SERVICE_URL}/api/v1/products/`, {
                    params: {
                        page: 1,
                        limit: 10
                    }
                });

                // 根據創建日期排序，最新的在前面
                const sortedProducts = response.data.data.products.sort((a, b) =>
                    new Date(b.created_at) - new Date(a.created_at)
                );

                setProducts(sortedProducts);
                // 設置輪播商品（取前5個）
                setFeaturedProducts(sortedProducts.slice(0, 5));
            } catch (error) {
                console.error('Error fetching products:', error);
                if (error.response?.status === 401) {
                    navigate('/login');
                }
            } finally {
                setLoading(false);

                // 頁面加載完成後再加載收藏狀態
                fetchWishlist();
            }
        };

        fetchProducts();
    }, [navigate, authAxios, fetchWishlist]);

    // 添加到收藏
    const handleToggleFavorite = useCallback(async (e, productId) => {
        e.stopPropagation(); // 阻止事件冒泡到卡片點擊事件

        try {
            if (favorites[productId]) {
                // 從收藏中移除
                await authAxios.delete(`${USER_SERVICE_URL}/api/v1/wishlist/${productId}`);
                setFavorites(prev => ({ ...prev, [productId]: false }));
                showAlert('已從收藏清單移除');
            } else {
                // 添加到收藏
                await authAxios.post(`${USER_SERVICE_URL}/api/v1/wishlist/`, {
                    productId: productId
                });
                setFavorites(prev => ({ ...prev, [productId]: true }));
                showAlert('已加入收藏');
            }
        } catch (error) {
            console.error('Error toggling favorite:', error);
            showAlert('操作失敗，請稍後再試', 'error');
        }
    }, [authAxios, favorites, showAlert]);

    // 添加到購物車
    const handleAddToCart = useCallback(async (e, productId) => {
        e.stopPropagation(); // 阻止事件冒泡到卡片點擊事件

        try {
            await authAxios.post(`${CART_SERVICE_URL}/api/v1/cart/items`, {
                productId: productId,
                quantity: 1
            });
            showAlert('已加入購物車');
        } catch (error) {
            console.error('Error adding to cart:', error);
            showAlert('加入購物車失敗', 'error');
        }
    }, [authAxios, showAlert]);

    // 自動輪播
    useEffect(() => {
        if (featuredProducts.length === 0) return;

        const timer = setInterval(() => {
            setCurrentProductIndex((prevIndex) =>
                prevIndex === featuredProducts.length - 1 ? 0 : prevIndex + 1
            );
        }, 5000); // 每5秒切換一次

        return () => clearInterval(timer);
    }, [featuredProducts]);

    const currentProduct = featuredProducts[currentProductIndex];

    // 導航處理函數，使用useCallback優化性能
    const handleNavigate = useCallback((path) => {
        navigate(path);
    }, [navigate]);

    // 橫幅骨架屏
    const BannerSkeleton = () => (
        <Paper
            sx={{
                position: 'relative',
                height: 500,
                mb: 6,
                borderRadius: '16px',
                overflow: 'hidden',
            }}
        >
            <Skeleton variant="rectangular" height="100%" animation="wave" />
        </Paper>
    );

    return (
        <Container maxWidth="xl">
            {/* 動態橫幅廣告 - 使用條件渲染顯示載入狀態 */}
            {loading ? <BannerSkeleton /> : currentProduct && (
                <Paper
                    sx={{
                        position: 'relative',
                        color: '#fff',
                        mb: 6,
                        backgroundSize: 'cover',
                        backgroundRepeat: 'no-repeat',
                        backgroundPosition: 'center',
                        backgroundImage: `url(${currentProduct.images && currentProduct.images[0]?.url ? currentProduct.images[0].url : '/images/default-banner.jpg'})`,
                        height: 500,
                        display: 'flex',
                        alignItems: 'center',
                        borderRadius: '16px',
                        overflow: 'hidden',
                        transition: 'background-image 0.5s ease-in-out',
                    }}
                >
                    <Box
                        sx={{
                            position: 'absolute',
                            top: 0,
                            bottom: 0,
                            right: 0,
                            left: 0,
                            background: 'linear-gradient(to right, rgba(0,0,0,0.8) 0%, rgba(0,0,0,0.4) 100%)',
                        }}
                    />
                    <Container maxWidth="lg">
                        <Box
                            sx={{
                                position: 'relative',
                                p: { xs: 3, md: 6 },
                                pr: { md: 0 },
                            }}
                        >
                            <Typography
                                component="h1"
                                variant="h2"
                                color="inherit"
                                gutterBottom
                                sx={{
                                    fontWeight: 700,
                                    textShadow: '0 2px 4px rgba(0,0,0,0.3)',
                                }}
                            >
                                {currentProduct.name}
                            </Typography>
                            <Typography
                                variant="h5"
                                color="inherit"
                                paragraph
                                sx={{
                                    mb: 4,
                                    textShadow: '0 1px 2px rgba(0,0,0,0.3)',
                                }}
                            >
                                {currentProduct.description}
                            </Typography>
                            <Button
                                variant="contained"
                                size="large"
                                onClick={() => navigate(`/store/products/${currentProduct.id}`)}
                                sx={{
                                    px: 4,
                                    py: 1.5,
                                    fontSize: '1.1rem',
                                    backgroundColor: 'white',
                                    color: '#1e293b',
                                    '&:hover': {
                                        backgroundColor: 'rgba(255,255,255,0.9)',
                                    },
                                }}
                            >
                                立即選購
                            </Button>
                        </Box>
                    </Container>
                </Paper>
            )}

            {/* 輪播指示器 */}
            {!loading && featuredProducts.length > 0 && (
                <Box
                    sx={{
                        display: 'flex',
                        justifyContent: 'center',
                        gap: 1,
                        mt: -4,
                        mb: 4,
                    }}
                >
                    {featuredProducts.map((_, index) => (
                        <Box
                            key={index}
                            onClick={() => setCurrentProductIndex(index)}
                            sx={{
                                width: 12,
                                height: 12,
                                borderRadius: '50%',
                                backgroundColor: index === currentProductIndex ? 'primary.main' : 'grey.300',
                                cursor: 'pointer',
                                transition: 'all 0.3s ease',
                                '&:hover': {
                                    transform: 'scale(1.2)',
                                },
                            }}
                        />
                    ))}
                </Box>
            )}

            {/* 商品列表 */}
            <Box sx={{ mb: 8 }}>
                <Typography
                    variant="h4"
                    gutterBottom
                    sx={{
                        mb: 4,
                        fontWeight: 700,
                    }}
                >
                    熱門商品
                </Typography>
                <Grid container spacing={4}>
                    {loading ? (
                        // 顯示骨架屏
                        Array.from(new Array(8)).map((_, index) => (
                            <Grid item key={index} xs={12} sm={6} md={4} lg={3}>
                                <ProductSkeleton />
                            </Grid>
                        ))
                    ) : (
                        products.map((product) => (
                            <Grid item key={product.id} xs={12} sm={6} md={4} lg={3}>
                                <ProductCard
                                    product={product}
                                    isFavorite={favorites[product.id]}
                                    onToggleFavorite={handleToggleFavorite}
                                    onAddToCart={handleAddToCart}
                                    onNavigate={handleNavigate}
                                />
                            </Grid>
                        ))
                    )}
                </Grid>
            </Box>

            {/* 提示信息 */}
            <Snackbar
                open={alert.open}
                autoHideDuration={3000}
                onClose={handleCloseAlert}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            >
                <Alert onClose={handleCloseAlert} severity={alert.severity}>
                    {alert.message}
                </Alert>
            </Snackbar>
        </Container>
    );
};

export default React.memo(StorePage); 