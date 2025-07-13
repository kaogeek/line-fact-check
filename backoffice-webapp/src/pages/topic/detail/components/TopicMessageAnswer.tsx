import { TYH3 } from '@/components/Typography';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

interface TopicMessageAnswerProps {
  onClickHistory: () => void;
}

const FormSchema = z.object({
  type: z.enum(['real', 'fake'], {
    required_error: 'You need to select a answser type.',
  }),
  answer: z.string(),
});

export default function TopicMessageAnswer({ onClickHistory }: TopicMessageAnswerProps) {
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
  });

  function onSubmit(data: z.infer<typeof FormSchema>) {
    console.log(data);
  }

  return (
    <>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col gap-2">
            <div className="flex gap-2">
              <TYH3 className="flex-1">Answer</TYH3>
              <Button variant="outline" type="button" onClick={onClickHistory}>
                History
              </Button>
              <Button variant="default" type="submit">
                Save
              </Button>
            </div>
            <FormField
              control={form.control}
              name="answer"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Textarea
                      placeholder="Type your answer here."
                      rows={20}
                      value={field.value}
                      onChange={field.onChange}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            ></FormField>
            {/* TODO: find why this overflow cause Radio group */}
            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem className="space-y-3">
                  <FormLabel>Type</FormLabel>
                  <FormControl>
                    <RadioGroup onValueChange={field.onChange} defaultValue={field.value}>
                      <FormItem className="flex items-center gap-3">
                        <FormControl>
                          <RadioGroupItem value="real" />
                        </FormControl>
                        <FormLabel className="font-normal">Real</FormLabel>
                      </FormItem>
                      <FormItem className="flex items-center gap-3">
                        <FormControl>
                          <RadioGroupItem value="fake" />
                        </FormControl>
                        <FormLabel className="font-normal">Fake</FormLabel>
                      </FormItem>
                    </RadioGroup>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            ></FormField>
          </div>
        </form>
      </Form>
    </>
  );
}
