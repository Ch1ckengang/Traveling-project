import Spinner from './Spinner';

const variantClassMap = {
  primary: 'bg-primary-500 text-white hover:bg-primary-600',
  accent: 'bg-accent-500 text-white hover:bg-accent-600',
  outline: 'border-2 border-primary-500 text-primary-500 hover:bg-primary-50',
  ghost: 'text-primary-500 hover:bg-primary-50',
  danger: 'bg-red-500 text-white hover:bg-red-600'
};

const sizeClassMap = {
  sm: 'h-9 px-3 text-sm',
  md: 'h-11 px-4 text-sm',
  lg: 'h-12 px-5 text-base'
};

const Button = ({
  children,
  variant = 'primary',
  size = 'md',
  disabled = false,
  loading = false,
  onClick,
  type = 'button',
  fullWidth = false,
  className = '',
  ...rest
}) => {
  const isDisabled = disabled || loading;

  const classes = [
    'inline-flex items-center justify-center gap-2 rounded-button font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-300 disabled:cursor-not-allowed disabled:opacity-60',
    variantClassMap[variant] || variantClassMap.primary,
    sizeClassMap[size] || sizeClassMap.md,
    fullWidth ? 'w-full' : '',
    className
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <button type={type} className={classes} disabled={isDisabled} onClick={onClick} {...rest}>
      {loading && <Spinner size='sm' color='text-current' />}
      {children}
    </button>
  );
};

export default Button;
