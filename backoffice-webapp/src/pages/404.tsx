import { Home } from 'lucide-react';
import NoDataState from '../components/state/NoDataState';
import { Button } from '../components/ui/button';

export default function NotFoundPage() {
  return (
    <div className="flex items-center justify-center h-screen">
      <NoDataState
        title="404"
        msg="The page you are looking for doesn't exist or has been moved."
        action={
          <Button variant="default">
            <Home />
            <a href="/">Go back home</a>
          </Button>
        }
      />
    </div>
  );
}
