import React from 'react';
import {
    Box,
    Container,
    Grid,
    Paper,
    Typography,
    List,
    ListItem,
    ListItemText,
    Divider,
} from '@mui/material';
import {
    TrendingUp as TrendingUpIcon,
    ShoppingCart as OrderIcon,
    People as CustomerIcon,
    Inventory as ProductIcon,
} from '@mui/icons-material';

const StatCard = ({ title, value, icon: Icon, color }) => (
    <Paper
        sx={{
            p: 3,
            display: 'flex',
            alignItems: 'center',
            gap: 2,
            height: '100%',
        }}
    >
        <Box
            sx={{
                backgroundColor: `${color}.lighter`,
                borderRadius: 2,
                p: 1,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
            }}
        >
            <Icon sx={{ fontSize: 40, color: `${color}.main` }} />
        </Box>
        <Box>
            <Typography variant="h6" color="text.secondary">
                {title}
            </Typography>
            <Typography variant="h4">
                {value}
            </Typography>
        </Box>
    </Paper>
);

const RecentOrdersList = () => (
    <Paper sx={{ p: 2, height: '100%' }}>
        <Typography variant="h6" gutterBottom>
            最近訂單
        </Typography>
        <List>
            <ListItem>
                <ListItemText
                    primary="訂單 #12345"
                    secondary="2024/01/15 10:30"
                />
                <Typography variant="body2" color="primary">
                    $1,299
                </Typography>
            </ListItem>
            <Divider />
            <ListItem>
                <ListItemText
                    primary="訂單 #12344"
                    secondary="2024/01/15 09:45"
                />
                <Typography variant="body2" color="primary">
                    $2,499
                </Typography>
            </ListItem>
            <Divider />
            <ListItem>
                <ListItemText
                    primary="訂單 #12343"
                    secondary="2024/01/15 09:15"
                />
                <Typography variant="body2" color="primary">
                    $899
                </Typography>
            </ListItem>
        </List>
    </Paper>
);

const TopProductsList = () => (
    <Paper sx={{ p: 2, height: '100%' }}>
        <Typography variant="h6" gutterBottom>
            熱銷商品
        </Typography>
        <List>
            <ListItem>
                <ListItemText
                    primary="iPhone 15 Pro"
                    secondary="銷售量：156"
                />
                <Typography variant="body2" color="primary">
                    $39,900
                </Typography>
            </ListItem>
            <Divider />
            <ListItem>
                <ListItemText
                    primary="MacBook Air M2"
                    secondary="銷售量：89"
                />
                <Typography variant="body2" color="primary">
                    $42,900
                </Typography>
            </ListItem>
            <Divider />
            <ListItem>
                <ListItemText
                    primary="AirPods Pro"
                    secondary="銷售量：234"
                />
                <Typography variant="body2" color="primary">
                    $7,990
                </Typography>
            </ListItem>
        </List>
    </Paper>
);

const Dashboard = () => {
    return (
        <Container maxWidth="lg">
            <Typography variant="h4" sx={{ mb: 5 }}>
                歡迎回來
            </Typography>

            <Grid container spacing={3}>
                {/* 統計卡片 */}
                <Grid item xs={12} sm={6} md={3}>
                    <StatCard
                        title="本月營收"
                        value="$158,000"
                        icon={TrendingUpIcon}
                        color="success"
                    />
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <StatCard
                        title="訂單數量"
                        value="48"
                        icon={OrderIcon}
                        color="info"
                    />
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <StatCard
                        title="會員數量"
                        value="1,257"
                        icon={CustomerIcon}
                        color="warning"
                    />
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                    <StatCard
                        title="商品數量"
                        value="89"
                        icon={ProductIcon}
                        color="error"
                    />
                </Grid>

                {/* 最近訂單和熱銷商品 */}
                <Grid item xs={12} md={6}>
                    <RecentOrdersList />
                </Grid>
                <Grid item xs={12} md={6}>
                    <TopProductsList />
                </Grid>
            </Grid>
        </Container>
    );
};

export default Dashboard; 