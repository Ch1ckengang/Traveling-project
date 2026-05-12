const variantClassMap = {
  success: 'bg-green-100 text-green-700 border border-green-200',
  warning: 'bg-amber-100 text-amber-700 border border-amber-200',
  danger: 'bg-red-100 text-red-700 border border-red-200',
  info: 'bg-blue-100 text-blue-700 border border-blue-200',
  default: 'bg-slate-100 text-slate-700 border border-slate-200'
};

const Badge = ({ children, variant = 'default', className = '' }) => {
  const classes = [
    'inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium',
    variantClassMap[variant] || variantClassMap.default,
    className
  ]
    .filter(Boolean)
    .join(' ');

  return <span className={classes}>{children}</span>;
};

export default Badge;
