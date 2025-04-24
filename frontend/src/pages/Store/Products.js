import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
    Container,
    Grid,
    Paper,
    FormControl,
    Select,
    MenuItem,
    Typography,
    Box,
    Stack,
    Card,
    CardMedia,
    CardContent,
    IconButton,
    Pagination,
    CircularProgress,
    Alert,
    Slider,
    Drawer,
    useMediaQuery,
    useTheme,
    Snackbar,
} from '@mui/material';
import {
    Sort as SortIcon,
    ShoppingCart as CartIcon,
    Favorite as FavoriteIcon,
    FavoriteBorder as FavoriteBorderIcon,
    Close as CloseIcon,
} from '@mui/icons-material';
import { createAuthAxios } from '../../utils/auth';
import axios from 'axios';

const ITEMS_PER_PAGE = 12;
const PRICE_MARKS = [
    { value: 0, label: '$0' },
    { value: 1000, label: '$1000' },
    { value: 5000, label: '$5000' },
    { value: 10000, label: '$10000' },
];

const Products = () => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
    const location = useLocation();
    const navigate = useNavigate();
    const [products, setProducts] = useState([]);
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedCategory, setSelectedCategory] = useState('all');
    const [priceRange, setPriceRange] = useState([0, 10000]);
    const [sortBy, setSortBy] = useState('default');
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [filterDrawerOpen, setFilterDrawerOpen] = useState(false);
    const [wishlist, setWishlist] = useState([]);
    const [snackbar, setSnackbar] = useState({
        open: false,
        message: '',
        severity: 'success'
    });
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    const authAxios = useMemo(() => createAuthAxios(navigate), [navigate]);

    // 檢查用戶是否已登入
    useEffect(() => {
        const userToken = localStorage.getItem('userToken');
        setIsLoggedIn(!!userToken);
    }, []);

    // 從 URL 獲取參數
    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const q = params.get('q');
        const cat = params.get('category');
        const sort = params.get('sort');
        const pageNum = parseInt(params.get('page')) || 1;
        const minPrice = parseInt(params.get('minPrice'));
        const maxPrice = parseInt(params.get('maxPrice'));

        if (q) setSearchTerm(q);
        if (cat) setSelectedCategory(cat);
        if (sort) setSortBy(sort);
        if (pageNum) setPage(pageNum);
        if (minPrice !== undefined && maxPrice !== undefined) {
            setPriceRange([minPrice, maxPrice]);
        }
    }, [location.search]);

    // 使用 useCallback 包裝 updateUrlParams
    const updateUrlParams = useCallback(() => {
        const params = new URLSearchParams();
        if (searchTerm) params.set('q', searchTerm);
        if (selectedCategory !== 'all') params.set('category', selectedCategory);
        if (sortBy !== 'default') params.set('sort', sortBy);
        if (page > 1) params.set('page', page.toString());
        if (priceRange[0] > 0 || priceRange[1] < 10000) {
            params.set('minPrice', priceRange[0].toString());
            params.set('maxPrice', priceRange[1].toString());
        }

        navigate({
            pathname: '/store/products',
            search: params.toString()
        });
    }, [navigate, searchTerm, selectedCategory, sortBy, page, priceRange]);

    // 使用 useCallback 包裝 fetchProducts，改為使用普通 axios
    const fetchProducts = useCallback(async () => {
        setLoading(true);
        try {
            const params = new URLSearchParams();
            params.set('page', page.toString());
            params.set('limit', ITEMS_PER_PAGE.toString());
            if (searchTerm) params.set('query', searchTerm);
            if (selectedCategory !== 'all') params.set('category', selectedCategory);
            if (sortBy !== 'default') params.set('sort', sortBy);

            // 使用普通 axios 獲取產品資訊
            const response = await axios.get(process.env.REACT_APP_PRODUCT_SERVICE_URL + `/api/v1/products/?${params.toString()}`);

            if (response.data.success) {
                setProducts(response.data.data.products);
                setTotal(response.data.data.total);
                setError('');
            } else {
                setError(response.data.message || '獲取商品失敗');
            }
        } catch (error) {
            setError(error.response?.data?.message || '無法載入商品列表');
        } finally {
            setLoading(false);
        }
    }, [page, searchTerm, selectedCategory, sortBy]);

    // 新增獲取類別的函數 - 同樣使用普通 axios
    const fetchCategories = useCallback(async () => {
        try {
            const response = await axios.get(process.env.REACT_APP_PRODUCT_SERVICE_URL + '/api/v1/categories/');

            // 檢查回應是否為陣列
            if (Array.isArray(response.data)) {
                setCategories(response.data);
            } else {
                console.error('獲取類別失敗: 回應格式不正確');
                setCategories([]);
            }
        } catch (error) {
            console.error('獲取類別失敗:', error);
            setCategories([]);
        }
    }, []);

    // 在組件載入時獲取類別
    useEffect(() => {
        fetchCategories();
    }, [fetchCategories]);

    // 獲取收藏列表，只有在用戶登入時執行
    const fetchWishlist = useCallback(async () => {
        if (!isLoggedIn) return;

        try {
            const response = await authAxios.get(process.env.REACT_APP_USER_SERVICE_URL + `/api/v1/wishlist/ids`);
            if (response.data.success) {
                setWishlist(response.data.data.wishlistIds || []);
            }
        } catch (error) {
            console.error('獲取收藏列表失敗:', error);
            // 不需要跳轉到登入頁面
        }
    }, [authAxios, isLoggedIn]);

    // 在組件載入時獲取收藏列表
    useEffect(() => {
        fetchWishlist();
    }, [fetchWishlist]);

    // 監聽參數變化
    useEffect(() => {
        fetchProducts();
        updateUrlParams();
    }, [fetchProducts, updateUrlParams]);

    const handlePriceChange = (event, newValue) => {
        setPriceRange(newValue);
        setPage(1);
    };

    const handleSortChange = (event) => {
        setSortBy(event.target.value);
        setPage(1);
    };

    const handlePageChange = (event, value) => {
        setPage(value);
        window.scrollTo(0, 0);
    };

    // 判斷商品是否已收藏
    const isProductInWishlist = (productId) => {
        return wishlist.includes(productId);
    };

    // 添加或移除收藏 - 檢查登入狀態
    const handleToggleWishlist = async (productId, event) => {
        event.stopPropagation();

        // 檢查是否已登入
        if (!isLoggedIn) {
            setSnackbar({
                open: true,
                message: '請先登入才能使用收藏功能',
                severity: 'info'
            });
            return;
        }

        try {
            if (isProductInWishlist(productId)) {
                // 從收藏中移除
                const response = await authAxios.delete(process.env.REACT_APP_USER_SERVICE_URL + `/api/v1/wishlist/${productId}`);
                if (response.data.success) {
                    setWishlist(wishlist.filter(id => id !== productId));
                    setSnackbar({
                        open: true,
                        message: '已從收藏清單移除',
                        severity: 'success'
                    });
                } else {
                    setSnackbar({
                        open: true,
                        message: response.data.message || '移除收藏失敗',
                        severity: 'error'
                    });
                }
            } else {
                // 添加到收藏
                const response = await authAxios.post(process.env.REACT_APP_USER_SERVICE_URL + `/api/v1/wishlist`, {
                    productId
                });
                if (response.data.success) {
                    setWishlist([...wishlist, productId]);
                    setSnackbar({
                        open: true,
                        message: '已添加至收藏清單',
                        severity: 'success'
                    });
                } else {
                    setSnackbar({
                        open: true,
                        message: response.data.message || '添加收藏失敗',
                        severity: 'error'
                    });
                }
            }
        } catch (error) {
            setSnackbar({
                open: true,
                message: error.response?.data?.message || '操作失敗',
                severity: 'error'
            });
        }
    };

    // 添加到購物車 - 檢查登入狀態
    const handleAddToCart = async (productId, event) => {
        event.stopPropagation();

        // 檢查是否已登入
        if (!isLoggedIn) {
            setSnackbar({
                open: true,
                message: '請先登入才能使用購物車功能',
                severity: 'info'
            });
            return;
        }

        try {
            const response = await authAxios.post(process.env.REACT_APP_CART_SERVICE_URL + `/api/v1/cart/items`, {
                productId,
                quantity: 1
            });

            if (response.data.success) {
                setSnackbar({
                    open: true,
                    message: '已添加至購物車',
                    severity: 'success'
                });
            } else {
                setSnackbar({
                    open: true,
                    message: response.data.message || '添加購物車失敗',
                    severity: 'error'
                });
            }
        } catch (error) {
            setSnackbar({
                open: true,
                message: error.response?.data?.message || '添加購物車失敗',
                severity: 'error'
            });
        }
    };

    const handleCloseSnackbar = () => {
        setSnackbar({ ...snackbar, open: false });
    };

    const FilterDrawerContent = () => (
        <Box sx={{ width: 250, p: 2 }}>
            <Stack spacing={3}>
                <Box>
                    <Typography variant="subtitle1" gutterBottom>
                        商品分類
                    </Typography>
                    <FormControl fullWidth size="small">
                        <Select
                            value={selectedCategory}
                            onChange={(e) => setSelectedCategory(e.target.value)}
                        >
                            <MenuItem value="all">全部類別</MenuItem>
                            {categories && categories.length > 0 ? (
                                categories.map((category) => (
                                    <MenuItem key={category.id || category._id} value={category.id || category._id}>
                                        {category.name}
                                    </MenuItem>
                                ))
                            ) : (
                                <MenuItem disabled>無可用類別</MenuItem>
                            )}
                        </Select>
                    </FormControl>
                </Box>

                <Box>
                    <Typography variant="subtitle1" gutterBottom>
                        價格範圍
                    </Typography>
                    <Slider
                        value={priceRange}
                        onChange={handlePriceChange}
                        valueLabelDisplay="auto"
                        min={0}
                        max={10000}
                        marks={PRICE_MARKS}
                    />
                </Box>

                <Box>
                    <Typography variant="subtitle1" gutterBottom>
                        排序方式
                    </Typography>
                    <FormControl fullWidth size="small">
                        <Select
                            value={sortBy}
                            onChange={handleSortChange}
                        >
                            <MenuItem value="default">預設排序</MenuItem>
                            <MenuItem value="price_asc">價格由低到高</MenuItem>
                            <MenuItem value="price_desc">價格由高到低</MenuItem>
                            <MenuItem value="newest">最新上架</MenuItem>
                        </Select>
                    </FormControl>
                </Box>
            </Stack>
        </Box>
    );

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            {/* 搜尋和篩選區 */}
            <Paper elevation={1} sx={{ p: 2, mb: 3 }}>
                <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} md={6}>
                        <FormControl fullWidth size="small">
                            <Select
                                value={selectedCategory}
                                onChange={(e) => setSelectedCategory(e.target.value)}
                                displayEmpty
                            >
                                <MenuItem value="all">全部類別</MenuItem>
                                {categories && categories.length > 0 ? (
                                    categories.map((category) => (
                                        <MenuItem key={category.id || category._id} value={category.id || category._id}>
                                            {category.name}
                                        </MenuItem>
                                    ))
                                ) : (
                                    <MenuItem disabled>無可用類別</MenuItem>
                                )}
                            </Select>
                        </FormControl>
                    </Grid>

                    {!isMobile && (
                        <Grid item md={6}>
                            <Stack direction="row" spacing={2} justifyContent="flex-end">
                                <FormControl size="small" sx={{ minWidth: 120 }}>
                                    <Select
                                        value={sortBy}
                                        onChange={handleSortChange}
                                        displayEmpty
                                        startAdornment={<SortIcon sx={{ mr: 1 }} />}
                                    >
                                        <MenuItem value="default">預設排序</MenuItem>
                                        <MenuItem value="price_asc">價格由低到高</MenuItem>
                                        <MenuItem value="price_desc">價格由高到低</MenuItem>
                                        <MenuItem value="newest">最新上架</MenuItem>
                                    </Select>
                                </FormControl>

                            </Stack>
                        </Grid>
                    )}
                </Grid>
            </Paper>

            {/* 錯誤提示 */}
            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            {/* 商品列表 */}
            {loading ? (
                <Box display="flex" justifyContent="center" py={4}>
                    <CircularProgress />
                </Box>
            ) : products.length === 0 ? (
                <Alert severity="info">
                    沒有找到相關商品
                </Alert>
            ) : (
                <Grid container spacing={2}>
                    {products.map((product) => (
                        <Grid item key={product.id || product._id} xs={12} sm={6} md={4} lg={3}>
                            <Card
                                sx={{
                                    height: '100%',
                                    display: 'flex',
                                    flexDirection: 'column',
                                    transition: 'transform 0.2s',
                                    '&:hover': {
                                        transform: 'translateY(-4px)',
                                    },
                                }}
                            >
                                <CardMedia
                                    component="img"
                                    height="200"
                                    image={product.images?.[0]?.url || '/placeholder.png'}
                                    alt={product.name}
                                    sx={{ objectFit: 'contain', p: 2 }}
                                    onClick={() => navigate(`/store/products/${product.id || product._id}`)}
                                />
                                <CardContent sx={{ flexGrow: 1 }}>
                                    <Typography gutterBottom variant="h6" component="h2">
                                        {product.name}
                                    </Typography>
                                    <Typography variant="h6" color="primary" gutterBottom>
                                        NT$ {product.price.toLocaleString()}
                                    </Typography>
                                    <Stack direction="row" spacing={1} justifyContent="flex-end">
                                        <IconButton
                                            size="small"
                                            color="primary"
                                            onClick={(e) => handleToggleWishlist(product.id || product._id, e)}
                                        >
                                            {isProductInWishlist(product.id || product._id) ?
                                                <FavoriteIcon /> : <FavoriteBorderIcon />}
                                        </IconButton>
                                        <IconButton
                                            size="small"
                                            color="primary"
                                            onClick={(e) => handleAddToCart(product.id || product._id, e)}
                                        >
                                            <CartIcon />
                                        </IconButton>
                                    </Stack>
                                </CardContent>
                            </Card>
                        </Grid>
                    ))}
                </Grid>
            )}

            {/* 分頁 */}
            {total > ITEMS_PER_PAGE && (
                <Box sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
                    <Pagination
                        count={Math.ceil(total / ITEMS_PER_PAGE)}
                        page={page}
                        onChange={handlePageChange}
                        color="primary"
                        size={isMobile ? 'small' : 'medium'}
                    />
                </Box>
            )}

            {/* 篩選抽屜 */}
            <Drawer
                anchor={isMobile ? 'bottom' : 'right'}
                open={filterDrawerOpen}
                onClose={() => setFilterDrawerOpen(false)}
            >
                <Box sx={{ display: 'flex', justifyContent: 'flex-end', p: 1 }}>
                    <IconButton onClick={() => setFilterDrawerOpen(false)}>
                        <CloseIcon />
                    </IconButton>
                </Box>
                <FilterDrawerContent />
            </Drawer>

            {/* 提示消息 */}
            <Snackbar
                open={snackbar.open}
                autoHideDuration={3000}
                onClose={handleCloseSnackbar}
                message={snackbar.message}
            />
        </Container>
    );
};

export default Products; 