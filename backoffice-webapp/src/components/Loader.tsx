import { PulseLoader } from 'react-spinners';

interface LoaderProps {
  className?: string;
}

export default function Loader({ className }: LoaderProps) {
  return <PulseLoader className={className} color="var(--color-primary)" />;
}
