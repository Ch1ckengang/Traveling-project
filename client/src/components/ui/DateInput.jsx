import { useState, useEffect } from 'react';

/**
 * Custom Date Input với format dd/mm/yyyy
 * Hiển thị format Việt Nam nhưng value vẫn là yyyy-mm-dd (chuẩn ISO)
 */
const DateInput = ({ 
  value, 
  onChange, 
  min, 
  max, 
  required = false, 
  className = '',
  name = 'date'
}) => {
  const [displayValue, setDisplayValue] = useState('');

  // Convert yyyy-mm-dd → dd/mm/yyyy for display
  useEffect(() => {
    if (value) {
      const [year, month, day] = value.split('-');
      setDisplayValue(`${day}/${month}/${year}`);
    } else {
      setDisplayValue('');
    }
  }, [value]);

  // Parse dd/mm/yyyy → yyyy-mm-dd
  const handleChange = (e) => {
    const input = e.target.value;
    
    // Remove non-numeric characters except /
    const cleaned = input.replace(/[^\d/]/g, '');
    
    // Auto-add slashes
    let formatted = cleaned;
    if (cleaned.length >= 2 && !cleaned.includes('/')) {
      formatted = cleaned.slice(0, 2) + '/' + cleaned.slice(2);
    }
    if (cleaned.length >= 5 && cleaned.split('/').length === 2) {
      const parts = cleaned.split('/');
      formatted = parts[0] + '/' + parts[1].slice(0, 2) + '/' + parts[1].slice(2);
    }
    
    setDisplayValue(formatted);
    
    // Parse to yyyy-mm-dd if complete
    if (formatted.length === 10) {
      const [day, month, year] = formatted.split('/');
      if (day && month && year && year.length === 4) {
        const isoDate = `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
        
        // Validate date
        const date = new Date(isoDate);
        if (!isNaN(date.getTime())) {
          onChange({ target: { name, value: isoDate } });
        }
      }
    }
  };

  const handleBlur = () => {
    // Validate on blur
    if (displayValue && displayValue.length === 10) {
      const [day, month, year] = displayValue.split('/');
      const isoDate = `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
      const date = new Date(isoDate);
      
      if (isNaN(date.getTime())) {
        setDisplayValue('');
        onChange({ target: { name, value: '' } });
      }
    }
  };

  return (
    <div className="relative">
      <input
        type="text"
        value={displayValue}
        onChange={handleChange}
        onBlur={handleBlur}
        placeholder="dd/mm/yyyy"
        maxLength={10}
        required={required}
        className={`w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 ${className}`}
      />
      <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-gray-400">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
          <rect x="3" y="4" width="18" height="18" rx="2" stroke="currentColor" strokeWidth="2"/>
          <path d="M3 10h18M8 2v4M16 2v4" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
        </svg>
      </div>
    </div>
  );
};

export default DateInput;
