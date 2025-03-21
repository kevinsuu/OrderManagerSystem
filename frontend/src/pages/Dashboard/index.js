import React from 'react';
import {
    Grid,
    Paper,
    Typography,
    Box,
    Card,
    CardContent,
    IconButton,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Chip,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TrendingDownIcon from '@mui/icons-material/TrendingDown';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import PeopleIcon from '@mui/icons-material/People';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import InventoryIcon from '@mui/icons-material/Inventory';
import MoreVertIcon from '@mui/icons-material/MoreVert';

const StyledCard = styled(Card)(({ theme }) => ({
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    position: 'relative',
    '&:hover': {
        boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
    },
}));

const StatIcon = styled(Box)(({ theme, color }) => ({
    width: 48,
    height: 48,
    borderRadius: '50%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: theme.palette[color].light,
    color: theme.palette[color].main,
    marginBottom: theme.spacing(2),
}));

const recentOrders = [
    {
        id: '1',
        customer: '王小明',
        product: 'iPhone 15 Pro',
        amount: 39900,
        status: '已付款',
        date: '2024-03-21',
    },
    {
        id: '2',
        customer: '李小華',
        product: 'MacBook Air',
        amount: 42900,
        status: '處理中',
        date: '2024-03-21',
    },
    {
        id: '3',
        customer: '張小美',
        product: 'AirPods Pro',
        amount: 7990,
        status: '已出貨',
        date: '2024-03-20',
    },
    {
        id: '4',
        customer: '陳大寶',
        product: 'iPad Air',
        amount: 19900,
        status: '已完成',
        date: '2024-03-20',
    },
];

const Dashboard = () => {
    const getStatusColor = (status) => {
        switch (status) {
            case '已付款':
                return 'primary';
            case '處理中':
                return 'warning';
            case '已出貨':
                return 'info';
            case '已完成':
                return 'success';
            default:
                return 'default';
        }
    };

    return (
        <Box>
            <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
                儀表板
            </Typography>
            <Grid container spacing={3}>
                {/* 統計卡片 */}
                <Grid item xs={12} sm={6} md={3}>
                    <StyledCard>
                        <CardContent>
                            <StatIcon color="primary">
                                <ShoppingCartIcon />
                            </StatIcon>
                            <Typography variant="h6" color="textSecondary" gutterBottom>
                                今日訂單
                            </Typography>
                            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                <Typography variant="h4" component="div" sx={{ mr: 1 }}>
                                    24
                                </Typography>
                                <Chip
                                    icon={<TrendingUpIcon />}
                                    label="+12.5%"
                                    size="small"
                                    color="success"
                                />
                            </Box>
                            <Typography variant="body2" color="textSecondary">
                                較昨日
                            </Typography>
                        </CardContent>
                    </StyledCard>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <StyledCard>
                        <CardContent>
                            <StatIcon color="success">
                                <AttachMoneyIcon />
                            </StatIcon>
                            <Typography variant="h6" color="textSecondary" gutterBottom>
                                本月營收
                            </Typography>
                            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                <Typography variant="h4" component="div" sx={{ mr: 1 }}>
                                    $158,245
                                </Typography>
                                <Chip
                                    icon={<TrendingUpIcon />}
                                    label="+8.2%"
                                    size="small"
                                    color="success"
                                />
                            </Box>
                            <Typography variant="body2" color="textSecondary">
                                較上月
                            </Typography>
                        </CardContent>
                    </StyledCard>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <StyledCard>
                        <CardContent>
                            <StatIcon color="info">
                                <PeopleIcon />
                            </StatIcon>
                            <Typography variant="h6" color="textSecondary" gutterBottom>
                                活躍會員
                            </Typography>
                            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                <Typography variant="h4" component="div" sx={{ mr: 1 }}>
                                    1,285
                                </Typography>
                                <Chip
                                    icon={<TrendingUpIcon />}
                                    label="+5.3%"
                                    size="small"
                                    color="success"
                                />
                            </Box>
                            <Typography variant="body2" color="textSecondary">
                                較上月
                            </Typography>
                        </CardContent>
                    </StyledCard>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <StyledCard>
                        <CardContent>
                            <StatIcon color="warning">
                                <InventoryIcon />
                            </StatIcon>
                            <Typography variant="h6" color="textSecondary" gutterBottom>
                                庫存商品
                            </Typography>
                            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                <Typography variant="h4" component="div" sx={{ mr: 1 }}>
                                    324
                                </Typography>
                                <Chip
                                    icon={<TrendingDownIcon />}
                                    label="-2.8%"
                                    size="small"
                                    color="error"
                                />
                            </Box>
                            <Typography variant="body2" color="textSecondary">
                                較上月
                            </Typography>
                        </CardContent>
                    </StyledCard>
                </Grid>

                {/* 最近訂單 */}
                <Grid item xs={12}>
                    <Card>
                        <CardContent>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                                <Typography variant="h6">最近訂單</Typography>
                                <IconButton>
                                    <MoreVertIcon />
                                </IconButton>
                            </Box>
                            <TableContainer>
                                <Table>
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>訂單編號</TableCell>
                                            <TableCell>客戶名稱</TableCell>
                                            <TableCell>商品</TableCell>
                                            <TableCell align="right">金額</TableCell>
                                            <TableCell>狀態</TableCell>
                                            <TableCell>日期</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {recentOrders.map((order) => (
                                            <TableRow key={order.id} hover>
                                                <TableCell>#{order.id}</TableCell>
                                                <TableCell>{order.customer}</TableCell>
                                                <TableCell>{order.product}</TableCell>
                                                <TableCell align="right">
                                                    ${order.amount.toLocaleString()}
                                                </TableCell>
                                                <TableCell>
                                                    <Chip
                                                        label={order.status}
                                                        color={getStatusColor(order.status)}
                                                        size="small"
                                                    />
                                                </TableCell>
                                                <TableCell>{order.date}</TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </TableContainer>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};

export default Dashboard; 