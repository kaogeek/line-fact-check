import './App.css';
import { Button } from './components/ui/button';

// test linter

function App() {
  return (
    <>
      <div className="p-4 max-w-sm mx-auto bg-white rounded-xl shadow-md flex items-center space-x-4">
        <div>
          <div className="text-2xl font-medium text-black ">Tailwind & Shadcn</div>
          <p className="text-gray-500">Demonstration Component</p>
        </div>
        <Button variant="default">Click Me</Button>
      </div>
    </>
  );
}

export default App;
