import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
    Container,
    Grid,
    Card,
    CardContent,
    CardMedia,
    Typography,
    Box,
    Pagination,
    CircularProgress,
    Alert,
} from '@mui/material';
import axios from 'axios';

// 添加環境變數
const PRODUCT_SERVICE_URL = process.env.REACT_APP_PRODUCT_SERVICE_URL || 'https://ordermanagersystem-product-service.onrender.com';

const ProductSearchPage = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [total, setTotal] = useState(0);

    // 從 URL 獲取搜尋參數
    const searchParams = new URLSearchParams(location.search);
    const query = searchParams.get('query') || '';
    const page = parseInt(searchParams.get('page')) || 1;
    const limit = parseInt(searchParams.get('limit')) || 10;

    console.log('=== URL 參數 ===');
    console.log('location.search:', location.search);
    console.log('query:', query);
    console.log('page:', page);
    console.log('limit:', limit);

    useEffect(() => {
        const fetchSearchResults = async () => {
            try {
                setLoading(true);
                setError(null);
                const token = localStorage.getItem('userToken');
                const apiUrl = `${PRODUCT_SERVICE_URL}/api/v1/products/search?query=${encodeURIComponent(query)}&page=${page}&limit=${limit}`;

                console.log('=== API 請求 ===');
                console.log('API URL:', apiUrl);
                console.log('Token:', token ? '存在' : '不存在');

                const response = await axios.get(
                    apiUrl,
                    {
                        headers: {
                            'Authorization': `Bearer ${token}`
                        }
                    }
                );

                console.log('=== API 響應 ===');
                console.log('Response:', response.data);

                if (response.data.success) {
                    setProducts(response.data.data.products);
                    setTotal(response.data.data.total);
                } else {
                    setError(response.data.message || '搜尋失敗');
                }
            } catch (error) {
                console.error('=== API 錯誤 ===');
                console.error('Error:', error);
                console.error('Error response:', error.response);
                setError(error.response?.data?.message || '搜尋過程中發生錯誤');
            } finally {
                setLoading(false);
            }
        };

        if (query) {
            console.log('開始搜尋，query:', query);
            fetchSearchResults();
        } else {
            console.log('沒有搜尋關鍵字');
        }
    }, [query, page, limit]);

    const handlePageChange = (event, newPage) => {
        const newSearchParams = new URLSearchParams(location.search);
        newSearchParams.set('page', newPage);
        console.log("====newSearchParams", newSearchParams);

        navigate(`/store/products/search?${newSearchParams.toString()}`);
    };

    if (loading) {
        return (
            <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
                <CircularProgress />
            </Box>
        );
    }

    return (
        <Container sx={{ py: 4 }}>
            <Typography variant="h4" gutterBottom>
                搜尋結果：{query}
            </Typography>

            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            {!loading && products.length === 0 ? (
                <Alert severity="info">
                    沒有找到相關商品
                </Alert>
            ) : (
                <>
                    <Grid container spacing={4}>
                        {products.map((product) => (
                            <Grid item key={product.id} xs={12} sm={6} md={4} lg={3}>
                                <Card
                                    sx={{
                                        height: '100%',
                                        display: 'flex',
                                        flexDirection: 'column',
                                        cursor: 'pointer'
                                    }}
                                    onClick={() => navigate(`/store/products/${product.id}`)}
                                >
                                    <CardMedia
                                        component="img"
                                        sx={{
                                            height: 200,
                                            objectFit: 'contain',
                                            p: 2,
                                            bgcolor: 'background.paper'
                                        }}
                                        image={product.images?.[0]?.url || '/placeholder.png'}
                                        alt={product.name}
                                    />
                                    <CardContent sx={{ flexGrow: 1 }}>
                                        <Typography gutterBottom variant="h6" component="h2">
                                            {product.name}
                                        </Typography>
                                        <Typography color="text.secondary">
                                            NT$ {product.price.toLocaleString()}
                                        </Typography>
                                    </CardContent>
                                </Card>
                            </Grid>
                        ))}
                    </Grid>

                    {total > limit && (
                        <Box sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
                            <Pagination
                                count={Math.ceil(total / limit)}
                                page={page}
                                onChange={handlePageChange}
                                color="primary"
                            />
                        </Box>
                    )}
                </>
            )}
        </Container>
    );
};

export default ProductSearchPage; 