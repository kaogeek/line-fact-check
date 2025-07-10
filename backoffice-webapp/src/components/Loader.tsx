import { PulseLoader } from 'react-spinners';

interface LoaderProps {
  className?: string;
}

export default function Loader({ className }: LoaderProps) {
  // todo change color
  return <PulseLoader />;
}
