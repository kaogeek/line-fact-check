import { PulseLoader } from 'react-spinners';
import { useLoader } from '../hooks/useLoader';

interface LoaderProps {
  className?: string;
}

export default function Loader({ className }: LoaderProps) {
  const { isLoading } = useLoader();

  if (!isLoading) return null;

  return <PulseLoader className={className} color="var(--color-primary)" />;
}
