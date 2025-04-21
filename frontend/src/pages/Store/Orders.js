import React, { useState, useEffect, useMemo, useCallback } from 'react';
import {
    Container,
    Typography,
    Box,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Chip,
    CircularProgress,
    Alert,
    Button,
    Pagination,
    TextField,
    InputAdornment,
    IconButton,
    Tooltip,
} from '@mui/material';
import {
    Search as SearchIcon,
    Visibility as VisibilityIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { format } from 'date-fns';
import { createAuthAxios } from '../../utils/auth';

const ITEMS_PER_PAGE = 10;

// 在生產環境中使用環境變數
const API_URL = process.env.REACT_APP_ORDER_SERVICE_URL;

const Orders = () => {
    const navigate = useNavigate();
    const [allOrders, setAllOrders] = useState([]); // 存儲所有訂單
    const [displayedOrders, setDisplayedOrders] = useState([]); // 顯示的訂單（過濾後）
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [page, setPage] = useState(1);
    const [searchInput, setSearchInput] = useState('');

    const authAxios = useMemo(() => createAuthAxios(navigate), [navigate]);

    // 取得訂單資料
    const fetchOrders = useCallback(async () => {
        setLoading(true);
        try {
            const params = new URLSearchParams();
            params.set('page', '1'); // 先取得所有訂單，不分頁
            params.set('limit', '100'); // 設置較大的數量以獲取更多訂單

            // 實際 API 呼叫
            const response = await authAxios.get(`${API_URL}/api/v1/orders/?${params.toString()}`);

            console.log('訂單API回應:', response.data);

            if (response.data && Array.isArray(response.data.orders)) {
                setAllOrders(response.data.orders);
                setError('');
            } else {
                setError('獲取訂單數據格式錯誤');
                setAllOrders([]);
            }
        } catch (error) {
            console.error('獲取訂單失敗:', error);
            setError(error.response?.data?.message || '無法載入訂單列表');
            setAllOrders([]);
        } finally {
            setLoading(false);
        }
    }, [authAxios]);

    // 初始載入
    useEffect(() => {
        fetchOrders();
    }, [fetchOrders]);

    // 搜尋過濾訂單
    useEffect(() => {
        if (!searchInput.trim()) {
            // 如果搜尋輸入為空，顯示所有訂單（分頁）
            const startIndex = (page - 1) * ITEMS_PER_PAGE;
            const endIndex = startIndex + ITEMS_PER_PAGE;
            setDisplayedOrders(allOrders.slice(startIndex, endIndex));
            return;
        }

        // 搜尋字串轉小寫，以便不區分大小寫搜尋
        const searchTerm = searchInput.toLowerCase();

        // 過濾訂單
        const filteredOrders = allOrders.filter(order => {
            // 檢查訂單號
            const orderNumber = getOrderNumber(order).toLowerCase();
            if (orderNumber.includes(searchTerm)) return true;

            // 檢查商品名稱
            if (order.items && order.items.length > 0) {
                return order.items.some(item =>
                    (item.name || '').toLowerCase().includes(searchTerm)
                );
            }

            return false;
        });

        // 更新顯示的訂單
        const startIndex = (page - 1) * ITEMS_PER_PAGE;
        const endIndex = startIndex + ITEMS_PER_PAGE;
        setDisplayedOrders(filteredOrders.slice(startIndex, endIndex));
    }, [allOrders, searchInput, page]);

    // 頁面變更
    const handlePageChange = (event, value) => {
        setPage(value);
        window.scrollTo(0, 0);
    };

    // 搜尋輸入變更
    const handleSearchInputChange = (e) => {
        setSearchInput(e.target.value);
        setPage(1); // 重置頁碼
    };

    const getStatusColor = (status) => {
        switch (status) {
            case 'pending':
                return 'warning';
            case 'processing':
                return 'info';
            case 'shipped':
                return 'primary';
            case 'delivered':
                return 'success';
            case 'cancelled':
                return 'error';
            default:
                return 'default';
        }
    };

    const getStatusText = (status) => {
        switch (status) {
            case 'pending':
                return '待處理';
            case 'processing':
                return '處理中';
            case 'shipped':
                return '已出貨';
            case 'delivered':
                return '已送達';
            case 'cancelled':
                return '已取消';
            default:
                return status;
        }
    };

    const formatDate = (dateString) => {
        if (!dateString) return '無日期';
        try {
            return format(new Date(dateString), 'yyyy/MM/dd HH:mm');
        } catch (error) {
            console.error('日期格式化錯誤:', error);
            return dateString;
        }
    };

    // 生成訂單號 (如果API沒有提供orderNumber)
    const getOrderNumber = (order) => {
        if (order.orderNumber) return order.orderNumber;

        // 使用ID的前8位作為訂單號
        const shortId = order.id.substring(0, 8).toUpperCase();
        return `ORD-${shortId}`;
    };

    // 獲取訂單的主要商品名稱
    const getOrderItemsSummary = (order) => {
        if (!order.items || order.items.length === 0) {
            return '無商品資訊';
        }

        // 篩選出有名稱的商品
        const itemsWithName = order.items.filter(item => item.name);

        if (itemsWithName.length === 0) {
            // 如果沒有商品名稱，顯示商品數量
            return `${order.items.length} 件商品`;
        }

        if (itemsWithName.length === 1) {
            // 如果只有一種商品，顯示名稱和數量
            const item = itemsWithName[0];
            return `${item.name} x ${item.quantity}`;
        }

        // 如果有多種商品，顯示第一種商品名稱和總商品種類數
        return `${itemsWithName[0].name} 等 ${itemsWithName.length} 種商品`;
    };

    // 計算過濾後的總頁數
    const filteredTotalPages = useMemo(() => {
        if (!searchInput.trim()) {
            return Math.ceil(allOrders.length / ITEMS_PER_PAGE);
        }

        const searchTerm = searchInput.toLowerCase();
        const filteredCount = allOrders.filter(order => {
            const orderNumber = getOrderNumber(order).toLowerCase();
            if (orderNumber.includes(searchTerm)) return true;

            if (order.items && order.items.length > 0) {
                return order.items.some(item =>
                    (item.name || '').toLowerCase().includes(searchTerm)
                );
            }

            return false;
        }).length;

        return Math.ceil(filteredCount / ITEMS_PER_PAGE);
    }, [allOrders, searchInput]);

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Typography variant="h4" component="h1" gutterBottom>
                我的訂單
            </Typography>

            <Paper sx={{ p: 2, mb: 3 }}>
                <Box sx={{ mb: 2, display: 'flex' }}>
                    <TextField
                        fullWidth
                        size="small"
                        placeholder="搜尋訂單編號或商品名稱"
                        value={searchInput}
                        onChange={handleSearchInputChange}
                        InputProps={{
                            endAdornment: (
                                <InputAdornment position="end">
                                    <IconButton edge="end">
                                        <SearchIcon />
                                    </IconButton>
                                </InputAdornment>
                            ),
                        }}
                    />
                </Box>
            </Paper>

            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            {loading ? (
                <Box display="flex" justifyContent="center" py={4}>
                    <CircularProgress />
                </Box>
            ) : displayedOrders.length === 0 ? (
                <Alert severity="info">
                    {searchInput.trim() ? '沒有符合搜尋條件的訂單' : '您目前沒有任何訂單記錄'}
                </Alert>
            ) : (
                <TableContainer component={Paper}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>訂單編號</TableCell>
                                <TableCell>商品資訊</TableCell>
                                <TableCell>訂單日期</TableCell>
                                <TableCell>金額</TableCell>
                                <TableCell>狀態</TableCell>
                                <TableCell align="right">操作</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {displayedOrders.map((order) => (
                                <TableRow key={order.id}>
                                    <TableCell component="th" scope="row">
                                        {getOrderNumber(order)}
                                    </TableCell>
                                    <TableCell>
                                        <Tooltip title={order.items?.map(item => `${item.name || '未命名商品'} x ${item.quantity}`).join(', ') || '無商品資訊'}>
                                            <Typography noWrap sx={{ maxWidth: 200 }}>
                                                {getOrderItemsSummary(order)}
                                            </Typography>
                                        </Tooltip>
                                    </TableCell>
                                    <TableCell>{formatDate(order.createdAt)}</TableCell>
                                    <TableCell>NT$ {order.totalAmount.toLocaleString()}</TableCell>
                                    <TableCell>
                                        <Chip
                                            label={getStatusText(order.status)}
                                            color={getStatusColor(order.status)}
                                            size="small"
                                        />
                                    </TableCell>
                                    <TableCell align="right">
                                        <Button
                                            startIcon={<VisibilityIcon />}
                                            size="small"
                                            onClick={() => navigate(`/orders/${order.id}`)}
                                        >
                                            查看
                                        </Button>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}

            {filteredTotalPages > 1 && (
                <Box sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
                    <Pagination
                        count={filteredTotalPages}
                        page={page}
                        onChange={handlePageChange}
                        color="primary"
                    />
                </Box>
            )}
        </Container>
    );
};

export default Orders; 