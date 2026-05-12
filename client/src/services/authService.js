import axiosInstance from '../utils/axiosInstance';

export const login = async (email, password) => {
  const response = await axiosInstance.post('/login', {
    email,
    password
  });
  return response.data;
};

export const register = async (name, email, phone, password) => {
  const response = await axiosInstance.post('/register', {
    name,
    email,
    phone,
    password
  });
  return response.data;
};

export const logout = async () => {
  try {
    await axiosInstance.post('/logout');
  } catch {
    // Ignore errors, always clear local data
  }
};


export const refreshToken = async (token) => {
  const response = await axiosInstance.post('/token/refresh', {
    refreshToken: token
  });
  return response.data;
};

export const forgotPassword = async (email) => {
  const response = await axiosInstance.post('/password/forgot', {
    email
  });
  return response.data;
};

export const resetPassword = async (email, otpCode, newPassword) => {
  const response = await axiosInstance.post('/password/reset', {
    email,
    otp_code: otpCode,
    new_password: newPassword
  });
  return response.data;
};

export const sendOTP = async (email) => {
  const response = await axiosInstance.post('/otp/send', {
    email
  });
  return response.data;
};

export const verifyOTP = async (email, code) => {
  const response = await axiosInstance.post('/otp/verify', {
    email,
    code
  });
  return response.data;
};

export const getMe = async () => {
  const response = await axiosInstance.get('/users/me');
  return response.data;
};

export const updateProfile = async (userId, data) => {
  const response = await axiosInstance.put(`/users/${userId}`, data);
  return response.data;
};


export const changePassword = async (current, newPass) => {
  const response = await axiosInstance.put('/users/me/password', {
    current_password: current,
    new_password: newPass
  });
  return response.data;
};
