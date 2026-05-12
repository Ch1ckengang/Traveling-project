const Input = ({
  label,
  name,
  type = 'text',
  placeholder,
  value,
  onChange,
  error,
  hint,
  required = false,
  disabled = false,
  leftIcon,
  rightIcon,
  className = '',
  ...rest
}) => {
  const inputClasses = [
    'w-full rounded-button border bg-white py-2.5 text-sm outline-none transition-colors',
    leftIcon ? 'pl-10' : 'pl-3',
    rightIcon ? 'pr-10' : 'pr-3',
    error
      ? 'border-red-500 focus:border-red-500 focus:ring-2 focus:ring-red-200'
      : 'border-slate-300 focus:border-primary-500 focus:ring-2 focus:ring-primary-100',
    disabled ? 'cursor-not-allowed bg-slate-100 text-slate-500' : 'text-slate-900',
    className
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className='w-full'>
      {label && (
        <label htmlFor={name} className='mb-1.5 block text-sm font-medium text-slate-700'>
          {label}
          {required && <span className='ml-1 text-red-500'>*</span>}
        </label>
      )}

      <div className='relative'>
        {leftIcon && (
          <span className='pointer-events-none absolute inset-y-0 left-3 flex items-center text-slate-400'>
            {leftIcon}
          </span>
        )}

        <input
          id={name}
          name={name}
          type={type}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          disabled={disabled}
          required={required}
          className={inputClasses}
          {...rest}
        />

        {rightIcon && (
          <span className='absolute inset-y-0 right-3 flex items-center text-slate-400'>
            {rightIcon}
          </span>
        )}
      </div>

      {error ? (
        <p className='mt-1.5 text-xs text-red-600'>{error}</p>
      ) : (
        hint && <p className='mt-1.5 text-xs text-slate-500'>{hint}</p>
      )}
    </div>
  );
};

export default Input;
