import Button from './Button';

const EmptyState = ({
  title = 'Chua co du lieu',
  description = 'Hay thu lai voi bo loc khac hoac quay lai sau.',
  actionLabel,
  onAction,
  icon = '🧳'
}) => {
  return (
    <div className='mx-auto flex max-w-md flex-col items-center justify-center rounded-card border border-dashed border-slate-300 bg-white px-6 py-10 text-center'>
      <div className='mb-3 text-4xl' aria-hidden='true'>
        {icon}
      </div>

      <h3 className='text-lg font-semibold text-slate-900'>{title}</h3>
      <p className='mt-1 text-sm text-slate-500'>{description}</p>

      {actionLabel ? (
        <div className='mt-5'>
          <Button variant='outline' onClick={onAction}>
            {actionLabel}
          </Button>
        </div>
      ) : null}
    </div>
  );
};

export default EmptyState;
