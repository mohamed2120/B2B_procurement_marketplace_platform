import { ReactNode } from 'react';

interface CardProps {
  title?: string;
  children: ReactNode;
  className?: string;
  actions?: ReactNode;
}

export default function Card({ title, children, className = '', actions }: CardProps) {
  return (
    <div className={`bg-white rounded-lg shadow border ${className}`}>
      {(title || actions) && (
        <div className="px-6 py-4 border-b flex justify-between items-center">
          {title && <h3 className="text-lg font-semibold text-gray-900">{title}</h3>}
          {actions && <div>{actions}</div>}
        </div>
      )}
      <div className="px-6 py-4">{children}</div>
    </div>
  );
}
