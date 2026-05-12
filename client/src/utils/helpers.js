const STATUS_LABELS = {
  PENDING: 'Chờ xác nhận',
  CONFIRMED: 'Đã xác nhận',
  PAID: 'Đã thanh toán',
  UNPAID: 'Chưa thanh toán',
  CANCELLED: 'Đã hủy',
  COMPLETED: 'Hoàn thành',
  REFUNDED: 'Đã hoàn tiền'
};

const STATUS_COLORS = {
  PENDING: 'bg-amber-100 text-amber-700 border border-amber-200',
  CONFIRMED: 'bg-blue-100 text-blue-700 border border-blue-200',
  PAID: 'bg-emerald-100 text-emerald-700 border border-emerald-200',
  UNPAID: 'bg-orange-100 text-orange-700 border border-orange-200',
  CANCELLED: 'bg-rose-100 text-rose-700 border border-rose-200',
  COMPLETED: 'bg-teal-100 text-teal-700 border border-teal-200',
  REFUNDED: 'bg-slate-100 text-slate-700 border border-slate-200'
};

export const formatCurrency = (amount) => {
  const numericAmount = Number(amount || 0);
  return `${new Intl.NumberFormat('vi-VN').format(numericAmount)} \u20ab`;
};

export const formatDate = (isoString) => {
  if (!isoString) {
    return '';
  }

  const date = new Date(isoString);
  if (Number.isNaN(date.getTime())) {
    return '';
  }

  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();

  return `${day}/${month}/${year}`;
};

export const formatDatetime = (isoString) => {
  if (!isoString) {
    return '';
  }

  const date = new Date(isoString);
  if (Number.isNaN(date.getTime())) {
    return '';
  }

  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  return `${day}/${month}/${year} ${hours}:${minutes}`;
};

export const getStatusLabel = (status) => {
  const key = (status || '').toUpperCase();
  return STATUS_LABELS[key] || 'Không xác định';
};

export const getStatusColor = (status) => {
  const key = (status || '').toUpperCase();
  return STATUS_COLORS[key] || 'bg-gray-100 text-gray-700 border border-gray-200';
};

export const truncateText = (text, maxLength) => {
  const normalizedText = (text || '').toString();
  const limit = Number(maxLength) || 0;

  if (!limit || normalizedText.length <= limit) {
    return normalizedText;
  }

  return `${normalizedText.slice(0, limit).trim()}...`;
};

export const calculateNights = (days) => {
  const totalDays = Number(days) || 0;
  const totalNights = Math.max(totalDays - 1, 0);
  return `${totalDays} ngày ${totalNights} đêm`;
};
