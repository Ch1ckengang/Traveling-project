export const formatCurrency = (amount) => {
  const numeric = Number(amount);
  if (!Number.isFinite(numeric)) {
    return '0đ';
  }

  return `${new Intl.NumberFormat('vi-VN').format(numeric)}đ`;
};
