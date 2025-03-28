import React, { useState, useEffect } from 'react';
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
    CircularProgress,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import FavoriteIcon from '@mui/icons-material/Favorite';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

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

const LoadingContainer = styled(Box)({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '200px',
});

// 添加環境變數
const PRODUCT_SERVICE_URL = process.env.REACT_APP_PRODUCT_SERVICE_URL || 'https://ordermanagersystem-product-service.onrender.com';

const StorePage = () => {
    const navigate = useNavigate();
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [latestProduct, setLatestProduct] = useState(null);

    useEffect(() => {
        const fetchProducts = async () => {
            try {
                const token = localStorage.getItem('token');
                const response = await axios.get(`${PRODUCT_SERVICE_URL}/api/v1/products/`, {
                    params: {
                        page: 1,
                        limit: 10
                    },
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });

                // 根據創建日期排序，最新的在前面
                const sortedProducts = response.data.data.products.sort((a, b) =>
                    new Date(b.created_at) - new Date(a.created_at)
                );

                setProducts(sortedProducts);
                // 設置最新商品
                if (sortedProducts.length > 0) {
                    setLatestProduct(sortedProducts[0]);
                }
            } catch (error) {
                console.error('Error fetching products:', error);
                // 如果是 401 未授權錯誤，可以導向登入頁面
                if (error.response?.status === 401) {
                    navigate('/login');
                }
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, [navigate]);

    return (
        <Container maxWidth="xl">
            {/* 動態橫幅廣告 */}
            {latestProduct && (
                <Paper
                    sx={{
                        position: 'relative',
                        color: '#fff',
                        mb: 6,
                        backgroundSize: 'cover',
                        backgroundRepeat: 'no-repeat',
                        backgroundPosition: 'center',
                        backgroundImage: `url(${latestProduct.images[0]?.url || '/images/default-banner.jpg'})`,
                        height: 500,
                        display: 'flex',
                        alignItems: 'center',
                        borderRadius: '16px',
                        overflow: 'hidden',
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
                                {latestProduct.name}
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
                                {latestProduct.description}
                            </Typography>
                            <Button
                                variant="contained"
                                size="large"
                                onClick={() => navigate(`/store/product/${latestProduct.id}`)}
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
                {loading ? (
                    <LoadingContainer>
                        <CircularProgress />
                    </LoadingContainer>
                ) : (
                    <Grid container spacing={4}>
                        {products.map((product) => (
                            <Grid item key={product.id} xs={12} sm={6} md={4} lg={3}>
                                <StyledCard onClick={() => navigate(`/store/product/${product.id}`)}>
                                    <ProductImage
                                        image={product.images[0]?.url || 'default-image-url.jpg'}
                                        title={product.name}
                                    >
                                        {product.discount && (
                                            <DiscountBadge>
                                                {product.discount * 100}% OFF
                                            </DiscountBadge>
                                        )}
                                        <ProductActions>
                                            <IconButton size="small">
                                                <FavoriteIcon fontSize="small" />
                                            </IconButton>
                                            <IconButton size="small">
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
                                                ${product.price.toLocaleString()}
                                            </Typography>
                                        </Box>
                                    </CardContent>
                                </StyledCard>
                            </Grid>
                        ))}
                    </Grid>
                )}
            </Box>
        </Container>
    );
};

export default StorePage; 