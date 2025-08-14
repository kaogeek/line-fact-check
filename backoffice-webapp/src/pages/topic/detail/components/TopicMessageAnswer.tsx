import { TYH3 } from '@/components/Typography';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { useEffect } from 'react';
import { useGetTopicAnswerByTopicId } from '@/hooks/api/topicAnswer';
import { TopicAnswerType } from '@/lib/api/type/topic-answer';
import LoadingState from '@/components/state/LoadingState';
import ErrorState from '@/components/state/ErrorState';
import { useTranslation } from 'react-i18next';

interface TopicMessageAnswerProps {
  onClickHistory: () => void;
  topicId: string;
  onUpdateAnswer: (content: string) => Promise<void>;
}

const formSchema = z.object({
  type: z.enum([TopicAnswerType.REAL, TopicAnswerType.FAKE]),
  answer: z.string(),
});

export default function TopicMessageAnswer({ onClickHistory, topicId, onUpdateAnswer }: TopicMessageAnswerProps) {
  const { t } = useTranslation();
  const { isLoading, data: answer, error } = useGetTopicAnswerByTopicId(topicId);
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  useEffect(() => {
    if (answer) {
      form.reset({
        type: answer.type,
        answer: answer.text,
      });
    }
  }, [answer, form]);

  async function handleSubmit(data: z.infer<typeof formSchema>) {
    if (answer) {
      await onUpdateAnswer(data.answer);
    }
  }

  if (isLoading) {
    return <LoadingState />;
  }

  if (error) {
    console.log(error);
    return <ErrorState />;
  }

  return (
    <>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)}>
          <div className="flex flex-col gap-2">
            <div className="flex gap-2">
              <TYH3 className="flex-1">{t('topicMessageAnswer.answerLabel')}</TYH3>
              <Button variant="outline" type="button" onClick={onClickHistory}>
                {t('topicMessageAnswer.historyButton')}
              </Button>
              <Button variant="default" type="submit">
                {t('topicMessageAnswer.saveButton')}
              </Button>
            </div>
            <FormField
              control={form.control}
              name="answer"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Textarea
                      placeholder={t('topicMessageAnswer.answerPlaceholder')}
                      rows={20}
                      value={field.value}
                      onChange={field.onChange}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            ></FormField>
            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem className="space-y-3">
                  <FormLabel>{t('topicMessageAnswer.typeLabel')}</FormLabel>
                  <FormControl>
                    <RadioGroup onValueChange={field.onChange} value={field.value}>
                      <FormItem className="flex items-center gap-3">
                        <FormControl>
                          <RadioGroupItem value={TopicAnswerType.REAL} />
                        </FormControl>
                        <FormLabel className="font-normal">{t('topicMessageAnswer.realOption')}</FormLabel>
                      </FormItem>
                      <FormItem className="flex items-center gap-3">
                        <FormControl>
                          <RadioGroupItem value={TopicAnswerType.FAKE} />
                        </FormControl>
                        <FormLabel className="font-normal">{t('topicMessageAnswer.fakeOption')}</FormLabel>
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
