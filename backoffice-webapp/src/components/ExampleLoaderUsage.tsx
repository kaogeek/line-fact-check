import { useLoader } from '../hooks/loader';

export function ExampleLoaderUsage() {
  const { startLoading, stopLoading } = useLoader();

  const handleAction = async () => {
    startLoading();

    // Simulate async operation
    await new Promise((resolve) => setTimeout(resolve, 2000));

    stopLoading();
  };

  return (
    <div>
      <button onClick={handleAction}>Trigger Loader</button>
    </div>
  );
}
