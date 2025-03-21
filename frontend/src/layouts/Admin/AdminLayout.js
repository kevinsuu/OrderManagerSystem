import React from 'react';
import { Outlet } from 'react-router-dom';
import { Box } from '@mui/material';

const AdminLayout = () => {
    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            {/* 這裡可以添加管理後台的導航欄等組件 */}
            <Box component="main" sx={{ flexGrow: 1 }}>
                <Outlet />
            </Box>
        </Box>
    );
};

export default AdminLayout; 