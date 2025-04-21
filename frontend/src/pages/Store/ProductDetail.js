import React, { useState, useEffect, useCallback } from 'react';
import {
    Container,
    Grid,
    Box,
    Paper,
    Typography,
    Button,
    IconButton,
    Divider,
    Rating,
    TextField,
    Alert,
    Snackbar,
    CircularProgress,
    Breadcrumbs,
    Tabs,
    Tab,
    CardMedia,
} from '@mui/material';
import {
    Add as AddIcon,
    Remove as RemoveIcon,
    Favorite as FavoriteIcon,
    FavoriteBorder as FavoriteBorderIcon,
    ShoppingCart as CartIcon,
    ArrowBack as ArrowBackIcon,
    Inventory as InventoryIcon,
} from '@mui/icons-material';
import { Link, useParams, useNavigate } from 'react-router-dom';
import { styled } from '@mui/material/styles';
import { createAuthAxios } from '../../utils/auth';

// 服務URL常量
const PRODUCT_SERVICE_URL = process.env.REACT_APP_PRODUCT_SERVICE_URL || 'https://ordermanagersystem-product-service.onrender.com';
const CART_SERVICE_URL = process.env.REACT_APP_CART_SERVICE_URL || 'https://ordermanagersystem.onrender.com';

// 樣式組件
const ProductImage = styled(CardMedia)(({ theme }) => ({
    height: 500,
    borderRadius: theme.shape.borderRadius,
    backgroundColor: '#f8fafc',
    backgroundSize: 'contain',
}));

const ThumbnailImage = styled(Paper)(({ theme, selected }) => ({
    width: 80,
    height: 80,
    cursor: 'pointer',
    transition: 'all 0.2s',
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    border: selected ? `2px solid ${theme.palette.primary.main}` : 'none',
    '&:hover': {
        transform: 'scale(1.05)',
    },
}));

const StockIndicator = styled(Box)(({ theme, inStock }) => ({
    display: 'inline-flex',
    alignItems: 'center',
    padding: '4px 12px',
    borderRadius: 16,
    backgroundColor: inStock ? '#e8f5e9' : '#ffebee',
    color: inStock ? '#1b5e20' : '#c62828',
    fontSize: '0.875rem',
    fontWeight: 500,
    marginBottom: theme.spacing(2),
}));

const PriceDisplay = styled(Typography)(({ theme }) => ({
    fontWeight: 600,
    color: theme.palette.primary.main,
    marginBottom: theme.spacing(2),
}));

const TabPanel = ({ children, value, index, ...other }) => (
    <div
        role="tabpanel"
        hidden={value !== index}
        id={`product-tabpanel-${index}`}
        aria-labelledby={`product-tab-${index}`}
        {...other}
    >
        {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
);

const ProductDetail = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [product, setProduct] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [quantity, setQuantity] = useState(1);
    const [activeImage, setActiveImage] = useState(0);
    const [tabValue, setTabValue] = useState(0);
    const [alertOpen, setAlertOpen] = useState(false);
    const [alertMessage, setAlertMessage] = useState('');
    const [alertSeverity, setAlertSeverity] = useState('success');
    const [isFavorite, setIsFavorite] = useState(false);
    const [categoryInfo, setCategoryInfo] = useState(null);

    // 創建授權axios客戶端
    const authAxios = React.useMemo(() => createAuthAxios(navigate), [navigate]);

    // 成功消息顯示
    const showSuccessMessage = useCallback((message) => {
        setAlertMessage(message);
        setAlertSeverity('success');
        setAlertOpen(true);
    }, []);

    // 錯誤消息顯示
    const showErrorMessage = useCallback((message) => {
        setAlertMessage(message);
        setAlertSeverity('error');
        setAlertOpen(true);
    }, []);

    // 檢查商品是否在收藏列表中
    const checkIfFavorite = useCallback(async (productId) => {
        try {
            const response = await authAxios.get(`${process.env.REACT_APP_USER_SERVICE_URL}/api/v1/wishlist/`);
            if (response.data && response.data.success && response.data.data && response.data.data.wishlist) {
                const wishlistItems = response.data.data.wishlist;
                const isFav = wishlistItems.some(item => item.productId === productId);
                setIsFavorite(isFav);
            }
        } catch (err) {
            console.error('檢查收藏狀態失敗:', err);
        }
    }, [authAxios]);

    // 頁面載入時滾動到頂部
    useEffect(() => {
        window.scrollTo(0, 0);
    }, [id]);

    // 獲取商品詳情
    useEffect(() => {
        const fetchProductDetail = async () => {
            setLoading(true);
            try {
                const response = await authAxios.get(`${PRODUCT_SERVICE_URL}/api/v1/products/${id}`);

                // 檢查API響應格式
                if (response.data && response.data.data) {
                    setProduct(response.data.data);
                } else {
                    setProduct(response.data);
                }
            } catch (err) {
                console.error('獲取商品詳情失敗:', err);
                setError('無法加載商品詳情，請稍後再試');
                showErrorMessage('獲取商品詳情失敗');
            } finally {
                setLoading(false);
            }
        };

        fetchProductDetail();
    }, [id, authAxios, showErrorMessage]);

    // 檢查收藏狀態
    useEffect(() => {
        if (id) {
            checkIfFavorite(id);
        }
    }, [id, checkIfFavorite]);

    // 處理收藏切換
    const handleToggleFavorite = async () => {
        try {
            if (isFavorite) {
                await authAxios.delete(`${process.env.REACT_APP_USER_SERVICE_URL}/api/v1/wishlist/${id}`);
                showSuccessMessage('商品已從收藏列表移除');
            } else {
                await authAxios.post(`${process.env.REACT_APP_USER_SERVICE_URL}/api/v1/wishlist/`, {
                    productId: id,
                });
                showSuccessMessage('商品已加入收藏');
            }
            setIsFavorite(!isFavorite);
        } catch (err) {
            console.error('切換收藏狀態失敗:', err);
            showErrorMessage('操作失敗，請稍後再試');
        }
    };

    // 處理數量變更
    const handleQuantityChange = (event) => {
        const value = parseInt(event.target.value, 10);
        if (!isNaN(value)) {
            if (product && product.stock) {
                setQuantity(Math.max(1, Math.min(value, product.stock)));
            } else {
                setQuantity(Math.max(1, value));
            }
        }
    };

    // 增減數量
    const increaseQuantity = () => {
        if (product && product.stock && quantity < product.stock) {
            setQuantity(quantity + 1);
        } else if (!product || !product.stock) {
            setQuantity(quantity + 1);
        }
    };

    const decreaseQuantity = () => {
        if (quantity > 1) {
            setQuantity(quantity - 1);
        }
    };

    // 添加到購物車
    const addToCart = async () => {
        try {
            await authAxios.post(`${CART_SERVICE_URL}/api/v1/cart/items`, {
                productId: id,
                quantity: quantity
            });
            showSuccessMessage('商品已添加到購物車');
        } catch (err) {
            console.error('添加購物車失敗:', err);
            showErrorMessage('添加購物車失敗，請稍後再試');
        }
    };

    // 處理標籤切換
    const handleTabChange = (event, newValue) => {
        setTabValue(newValue);
    };

    // 返回上一頁
    const handleGoBack = () => {
        navigate(-1);
    };

    // 處理分類詳情獲取
    useEffect(() => {
        const fetchCategoryDetail = async () => {
            if (product && product.category && typeof product.category === 'string') {
                try {
                    console.log('正在獲取分類詳情...');
                    const response = await authAxios.get(`${PRODUCT_SERVICE_URL}/api/v1/categories/${product.category}`);
                    console.log('分類詳情:', response.data);
                    if (response.data) {
                        setCategoryInfo(response.data);
                    }
                } catch (err) {
                    console.error('獲取分類詳情失敗:', err);
                }
            }
        };

        if (product) {
            fetchCategoryDetail();
        }
    }, [product, authAxios]);

    // 更新渲染分類的函數
    const renderCategory = () => {
        if (!product.category) return null;

        let categoryName = '未知分類';

        // 處理不同的資料結構
        if (typeof product.category === 'object' && product.category.name) {
            categoryName = product.category.name;
        } else if (categoryInfo && categoryInfo.name) {
            // 使用從API獲取的分類資訊
            categoryName = categoryInfo.name;
        } else if (typeof product.category === 'string') {
            // 如果還沒獲取到分類詳情，顯示ID
            categoryName = product.category;
        }

        return (
            <Typography variant="subtitle2" color="text.secondary">
                分類: {categoryName}
            </Typography>
        );
    };

    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ my: 4 }}>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}>
                    <CircularProgress />
                </Box>
            </Container>
        );
    }

    if (error) {
        return (
            <Container maxWidth="lg" sx={{ my: 4 }}>
                <Alert severity="error">{error}</Alert>
            </Container>
        );
    }

    if (!product) {
        return (
            <Container maxWidth="lg" sx={{ my: 4 }}>
                <Alert severity="warning">找不到該商品</Alert>
            </Container>
        );
    }

    return (
        <Container maxWidth="lg" sx={{ my: 4 }}>
            {/* 麵包屑導航 */}
            <Box sx={{ mb: 4 }}>
                <Breadcrumbs separator="›" aria-label="breadcrumb">
                    <Link
                        to="/"
                        style={{ textDecoration: 'none', color: 'inherit' }}
                    >
                        首頁
                    </Link>
                    <Link
                        to="/store/products"
                        style={{ textDecoration: 'none', color: 'inherit' }}
                    >
                        商店
                    </Link>
                    <Typography color="text.primary">{product.name}</Typography>
                </Breadcrumbs>
                <Button
                    startIcon={<ArrowBackIcon />}
                    onClick={handleGoBack}
                    sx={{ mt: 2 }}
                >
                    返回
                </Button>
            </Box>

            <Grid container spacing={4}>
                {/* 商品圖片 */}
                <Grid item xs={12} md={6}>
                    <Box sx={{ mb: 2 }}>
                        <ProductImage
                            image={product.images && product.images.length > 0
                                ? product.images[activeImage]?.url || '/images/default-product.jpg'
                                : '/images/default-product.jpg'
                            }
                        />
                    </Box>
                    {/* 縮略圖 */}
                    <Box sx={{ display: 'flex', gap: 2, mt: 2, overflowX: 'auto', pb: 2 }}>
                        {product.images && product.images.map((image, index) => (
                            <ThumbnailImage
                                key={index}
                                selected={activeImage === index}
                                onClick={() => setActiveImage(index)}
                                style={{
                                    backgroundImage: `url(${image.url || '/images/default-product.jpg'})`
                                }}
                            />
                        ))}
                    </Box>
                </Grid>

                {/* 商品詳情 */}
                <Grid item xs={12} md={6}>
                    <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: 700 }}>
                        {product.name}
                    </Typography>

                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                        <Rating
                            value={product.rating || 4.5}
                            precision={0.5}
                            readOnly
                        />
                        <Typography variant="body2" color="text.secondary" sx={{ ml: 1 }}>
                            ({product.reviewCount || 0} 評價)
                        </Typography>
                    </Box>

                    <StockIndicator inStock={product.stock > 0}>
                        <InventoryIcon fontSize="small" sx={{ mr: 0.5 }} />
                        {product.stock > 0 ? '有庫存' : '缺貨中'}
                    </StockIndicator>

                    <PriceDisplay variant="h4">
                        ${product.price?.toLocaleString()}
                    </PriceDisplay>

                    <Typography variant="body1" sx={{ mb: 3 }}>
                        {product.description}
                    </Typography>

                    <Divider sx={{ mb: 3 }} />

                    {/* 數量選擇和加入購物車 */}
                    <Box sx={{ mb: 3 }}>
                        <Typography variant="subtitle1" sx={{ mb: 1 }}>
                            數量
                        </Typography>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <IconButton
                                onClick={decreaseQuantity}
                                disabled={quantity <= 1}
                                size="small"
                            >
                                <RemoveIcon />
                            </IconButton>
                            <TextField
                                value={quantity}
                                onChange={handleQuantityChange}
                                variant="outlined"
                                size="small"
                                inputProps={{ min: 1, max: product.stock || 99, style: { textAlign: 'center' } }}
                                sx={{ width: '80px' }}
                            />
                            <IconButton
                                onClick={increaseQuantity}
                                disabled={product.stock && quantity >= product.stock}
                                size="small"
                            >
                                <AddIcon />
                            </IconButton>
                        </Box>
                    </Box>

                    <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
                        <Button
                            variant="contained"
                            size="large"
                            startIcon={<CartIcon />}
                            disabled={product.stock <= 0}
                            onClick={addToCart}
                            fullWidth
                        >
                            加入購物車
                        </Button>
                        <IconButton
                            color={isFavorite ? "error" : "default"}
                            onClick={handleToggleFavorite}
                            sx={{ border: '1px solid', borderColor: 'divider' }}
                        >
                            {isFavorite ? <FavoriteIcon /> : <FavoriteBorderIcon />}
                        </IconButton>
                    </Box>

                    <Box sx={{ mb: 3 }}>
                        <Typography variant="subtitle2" color="text.secondary">
                            商品編號: {product.id}
                        </Typography>
                        {renderCategory()}
                    </Box>
                </Grid>
            </Grid>

            {/* 詳細內容標籤 */}
            <Box sx={{ mt: 6, mb: 4 }}>
                <Tabs
                    value={tabValue}
                    onChange={handleTabChange}
                    aria-label="product tabs"
                    sx={{ borderBottom: 1, borderColor: 'divider' }}
                >
                    <Tab label="商品詳情" />
                    <Tab label="規格參數" />
                    <Tab label="顧客評價" />
                </Tabs>

                <TabPanel value={tabValue} index={0}>
                    <Typography variant="body1">
                        {product.description}
                    </Typography>
                </TabPanel>

                <TabPanel value={tabValue} index={1}>
                    <Grid container spacing={2}>
                        {product.attributes && product.attributes.map((attr, index) => (
                            <Grid item xs={12} sm={6} key={index}>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between', py: 1 }}>
                                    <Typography variant="body2" color="text.secondary">
                                        {attr.name}:
                                    </Typography>
                                    <Typography variant="body2">
                                        {attr.value}
                                    </Typography>
                                </Box>
                                <Divider />
                            </Grid>
                        ))}
                        {(!product.attributes || product.attributes.length === 0) && (
                            <Grid item xs={12}>
                                <Typography variant="body2" color="text.secondary">
                                    暫無規格參數
                                </Typography>
                            </Grid>
                        )}
                    </Grid>
                </TabPanel>

                <TabPanel value={tabValue} index={2}>
                    <Typography variant="body2" color="text.secondary">
                        暫無評價
                    </Typography>
                </TabPanel>
            </Box>

            {/* 成功/錯誤提示 */}
            <Snackbar
                open={alertOpen}
                autoHideDuration={3000}
                onClose={() => setAlertOpen(false)}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            >
                <Alert onClose={() => setAlertOpen(false)} severity={alertSeverity}>
                    {alertMessage}
                </Alert>
            </Snackbar>
        </Container>
    );
};

export default ProductDetail; 