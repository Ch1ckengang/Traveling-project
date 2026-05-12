import { useEffect, useMemo, useState } from 'react';

const sizeClassMap = {
  sm: 'max-w-md',
  md: 'max-w-2xl',
  lg: 'max-w-4xl'
};

const Modal = ({ isOpen, onClose, title, children, size = 'md', footer }) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (!isOpen) {
      return;
    }

    const frame = requestAnimationFrame(() => {
      setIsVisible(true);
    });

    return () => cancelAnimationFrame(frame);
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen) {
      setIsVisible(false);
    }
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen) {
      return undefined;
    }

    const handleKeyDown = (event) => {
      if (event.key === 'Escape') {
        onClose?.();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  const panelClassName = useMemo(() => {
    return [
      'w-full rounded-card bg-white shadow-2xl transition-all duration-200',
      sizeClassMap[size] || sizeClassMap.md,
      isVisible ? 'translate-y-0 opacity-100' : 'translate-y-6 opacity-0'
    ]
      .filter(Boolean)
      .join(' ');
  }, [isVisible, size]);

  if (!isOpen) {
    return null;
  }

  return (
    <div
      className='fixed inset-0 z-50 flex items-end justify-center bg-black/40 p-4 md:items-center'
      onMouseDown={(event) => {
        if (event.target === event.currentTarget) {
          onClose?.();
        }
      }}
    >
      <div className={panelClassName} onMouseDown={(event) => event.stopPropagation()}>
        <header className='flex items-center justify-between border-b border-slate-200 px-4 py-3'>
          <h3 className='text-base font-semibold text-slate-900'>{title || 'Modal'}</h3>
          <button
            type='button'
            className='rounded-md p-1 text-slate-500 transition-colors hover:bg-slate-100 hover:text-slate-700'
            onClick={() => onClose?.()}
            aria-label='Close modal'
          >
            x
          </button>
        </header>

        <div className='max-h-[65vh] overflow-y-auto px-4 py-4'>{children}</div>

        {footer ? <footer className='border-t border-slate-200 px-4 py-3'>{footer}</footer> : null}
      </div>
    </div>
  );
};

export default Modal;
