import React from 'react';
import { useHealthCheck } from '../../utils/healthCheck';
import {
    Container,
    Paper,
    Typography,
    Grid,
    Box,
    CircularProgress
} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';

const SystemStatus = () => {
    const { serviceStatus } = useHealthCheck();

    return (
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Paper elevation={3} sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                    系統服務狀態
                </Typography>
                <Grid container spacing={3}>
                    {Object.entries(serviceStatus).map(([serviceName, isHealthy]) => (
                        <Grid item xs={12} sm={6} md={4} key={serviceName}>
                            <Box
                                sx={{
                                    p: 2,
                                    border: 1,
                                    borderColor: 'divider',
                                    borderRadius: 1,
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 2
                                }}
                            >
                                {isHealthy ? (
                                    <CheckCircleIcon color="success" />
                                ) : (
                                    <ErrorIcon color="error" />
                                )}
                                <Box>
                                    <Typography variant="subtitle1">
                                        {serviceName}
                                    </Typography>
                                    <Typography
                                        variant="body2"
                                        color={isHealthy ? 'success.main' : 'error.main'}
                                    >
                                        {isHealthy ? '正常運作中' : '服務異常'}
                                    </Typography>
                                </Box>
                            </Box>
                        </Grid>
                    ))}
                </Grid>
            </Paper>
        </Container>
    );
};

export default SystemStatus; 