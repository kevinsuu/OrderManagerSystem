import React from 'react';
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
} from '@mui/material';
import { styled } from '@mui/material/styles';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import FavoriteIcon from '@mui/icons-material/Favorite';
import { useNavigate } from 'react-router-dom';

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

const featuredProducts = [
    {
        id: 1,
        name: 'iPhone 15 Pro',
        description: '搭載 A17 Pro 晶片，4800 萬像素相機系統',
        price: 39900,
        image: 'https://store.storeimages.cdn-apple.com/8756/as-images.apple.com/is/iphone-15-pro-finish-select-202309-6-7inch_GEO_TW?wid=5120&hei=2880&fmt=p-jpg&qlt=80&.v=1693009279096',
        rating: 4.8,
        reviews: 128,
        discount: 0.1,
    },
    {
        id: 2,
        name: 'MacBook Air M2',
        description: '搭載 M2 晶片，超長續航力',
        price: 42900,
        image: 'https://store.storeimages.cdn-apple.com/8756/as-images.apple.com/is/macbook-air-midnight-config-20220606?wid=820&hei=498&fmt=jpeg&qlt=90&.v=1654122880566',
        rating: 4.9,
        reviews: 95,
    },
    {
        id: 3,
        name: 'AirPods Pro',
        description: '主動降噪，通透模式，完美音質',
        price: 7990,
        image: 'https://store.storeimages.cdn-apple.com/8756/as-images.apple.com/is/MQD83?wid=1144&hei=1144&fmt=jpeg&qlt=90&.v=1660803972361',
        rating: 4.7,
        reviews: 256,
        discount: 0.15,
    },
    {
        id: 4,
        name: 'iPad Air',
        description: 'M1 晶片，10.9 吋 Liquid Retina 顯示器',
        price: 19900,
        image: 'https://store.storeimages.cdn-apple.com/8756/as-images.apple.com/is/ipad-air-select-wifi-blue-202203?wid=940&hei=1112&fmt=png-alpha&.v=1645065732688',
        rating: 4.6,
        reviews: 184,
    },
];

const StorePage = () => {
    const navigate = useNavigate();

    return (
        <Container maxWidth="xl">
            {/* 橫幅廣告 */}
            <Paper
                sx={{
                    position: 'relative',
                    color: '#fff',
                    mb: 6,
                    backgroundSize: 'cover',
                    backgroundRepeat: 'no-repeat',
                    backgroundPosition: 'center',
                    backgroundImage: 'url(https://store.storeimages.cdn-apple.com/8756/as-images.apple.com/is/iphone-15-pro-model-unselect-gallery-2-202309_GEO_TW?wid=5120&hei=2880&fmt=p-jpg&qlt=80&.v=1693010535312)',
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
                        background: 'linear-gradient(to right, rgba(0,0,0,0.7) 0%, rgba(0,0,0,0.3) 100%)',
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
                            iPhone 15 Pro
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
                            鈦金屬設計，A17 Pro 晶片，
                            <br />
                            突破性相機系統
                        </Typography>
                        <Button
                            variant="contained"
                            size="large"
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
                    {featuredProducts.map((product) => (
                        <Grid item key={product.id} xs={12} sm={6} md={4} lg={3}>
                            <StyledCard onClick={() => navigate(`/store/product/${product.id}`)}>
                                <ProductImage
                                    image={product.image}
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
                                        <Rating value={product.rating} precision={0.1} readOnly size="small" />
                                        <Typography variant="body2" color="text.secondary" sx={{ ml: 1 }}>
                                            ({product.reviews})
                                        </Typography>
                                    </Box>
                                    <Box sx={{ display: 'flex', alignItems: 'baseline', gap: 1 }}>
                                        <Typography variant="h6" color="primary" sx={{ fontWeight: 600 }}>
                                            ${product.discount
                                                ? (product.price * (1 - product.discount)).toLocaleString()
                                                : product.price.toLocaleString()}
                                        </Typography>
                                        {product.discount && (
                                            <Typography
                                                variant="body2"
                                                color="text.secondary"
                                                sx={{ textDecoration: 'line-through' }}
                                            >
                                                ${product.price.toLocaleString()}
                                            </Typography>
                                        )}
                                    </Box>
                                </CardContent>
                            </StyledCard>
                        </Grid>
                    ))}
                </Grid>
            </Box>
        </Container>
    );
};

export default StorePage; 