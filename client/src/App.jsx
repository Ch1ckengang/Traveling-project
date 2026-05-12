import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import PublicLayout from './components/layouts/PublicLayout';
import AccountLayout from './components/layouts/AccountLayout';
import AdminLayout from './components/layouts/AdminLayout';
import HomePage from './pages/public/Home';
import TourListPage from './pages/public/TourList';
import TourDetailPage from './pages/public/TourDetail';
import SearchPage from './pages/public/Search';
import LoginPage from './pages/auth/Login';
import RegisterPage from './pages/auth/Register';
import ForgotPasswordPage from './pages/auth/ForgotPassword';
import ResetPasswordPage from './pages/auth/ResetPassword';
import OtpVerificationPage from './pages/auth/OtpVerification';
import ProfilePage from './pages/customer/Profile';
import BookingsPage from './pages/customer/Bookings';
import BookingDetailPage from './pages/customer/BookingDetail';
import WriteReviewPage from './pages/customer/WriteReview';
import AdminDashboardPage from './pages/admin/Dashboard';
import AdminToursPage from './pages/admin/Tours';
import AdminBookingsPage from './pages/admin/Bookings';
import AdminUsersPage from './pages/admin/Users';
import AdminCouponsPage from './pages/admin/Coupons';
import AdminReviewsPage from './pages/admin/Reviews';
import './App.css';

const AccountPasswordPage = () => <div>AccountPasswordPage</div>;
const AdminSchedulesPage = () => <div>AdminSchedulesPage</div>;
const AdminPaymentsPage = () => <div>AdminPaymentsPage</div>;
const AdminReportsPage = () => <div>AdminReportsPage</div>;

function AppRoutes() {
  return (
    <Routes>
      <Route element={<PublicLayout />}>
        <Route path='/' element={<HomePage />} />
        <Route path='/tours' element={<TourListPage />} />
        <Route path='/tours/:tourId' element={<TourDetailPage />} />
        <Route path='/search' element={<SearchPage />} />

        <Route path='/auth/login' element={<LoginPage />} />
        <Route path='/auth/register' element={<RegisterPage />} />
        <Route path='/auth/forgot-password' element={<ForgotPasswordPage />} />
        <Route path='/auth/reset-password' element={<ResetPasswordPage />} />
        <Route path='/auth/otp-verification' element={<OtpVerificationPage />} />

        <Route path='/login' element={<LoginPage />} />
        <Route path='/register' element={<RegisterPage />} />
        <Route path='/forgot-password' element={<ForgotPasswordPage />} />
        <Route path='/reset-password' element={<ResetPasswordPage />} />
        <Route path='/otp-verification' element={<OtpVerificationPage />} />
        <Route path='/reset-success' element={<ResetPasswordPage />} />
      </Route>

      <Route element={<AccountLayout />}>
        <Route path='/account/profile' element={<ProfilePage />} />
        <Route path='/account/bookings' element={<BookingsPage />} />
        <Route path='/account/bookings/:bookingId' element={<BookingDetailPage />} />
        <Route path='/account/reviews/write' element={<WriteReviewPage />} />
        <Route path='/account/password' element={<AccountPasswordPage />} />

        <Route path='/profile' element={<ProfilePage />} />
      </Route>

      <Route element={<AdminLayout />}>
        <Route path='/admin/dashboard' element={<AdminDashboardPage />} />
        <Route path='/admin/tours' element={<AdminToursPage />} />
        <Route path='/admin/schedules' element={<AdminSchedulesPage />} />
        <Route path='/admin/bookings' element={<AdminBookingsPage />} />
        <Route path='/admin/payments' element={<AdminPaymentsPage />} />
        <Route path='/admin/users' element={<AdminUsersPage />} />
        <Route path='/admin/reviews' element={<AdminReviewsPage />} />
        <Route path='/admin/coupons' element={<AdminCouponsPage />} />
        <Route path='/admin/reports' element={<AdminReportsPage />} />
      </Route>

      <Route path='*' element={<Navigate to='/' replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <AppRoutes />
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}
