const sizeClassMap = {
  sm: 'h-4 w-4 border-2',
  md: 'h-6 w-6 border-2',
  lg: 'h-8 w-8 border-[3px]'
};

const Spinner = ({ size = 'md', color = 'text-primary-500', className = '' }) => {
  const classes = [
    'inline-block animate-spin rounded-full border-current border-r-transparent',
    sizeClassMap[size] || sizeClassMap.md,
    color,
    className
  ]
    .filter(Boolean)
    .join(' ');

  return <span className={classes} aria-label='Loading' role='status' />;
};

export default Spinner;
