import React, { useEffect, useState, FormEvent } from 'react';
import { Form, Button } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import type { LangsType, FormDataType } from '@/common/interface';
import { useToast } from '@/hooks';
import { updateUserInterface } from '@/services';
import { localize } from '@/utils';
import { loggedUserInfoStore } from '@/stores';

const Index = () => {
  const { t } = useTranslation('translation', {
    keyPrefix: 'settings.interface',
  });
  const loggedUserInfo = loggedUserInfoStore.getState().user;
  const toast = useToast();
  const [langs, setLangs] = useState<LangsType[]>();
  const [formData, setFormData] = useState<FormDataType>({
    lang: {
      value: loggedUserInfo.language,
      isInvalid: false,
      errorMsg: '',
    },
  });

  const getLangs = async () => {
    const res: LangsType[] = await localize.loadLanguageOptions();
    setLangs(res);
  };

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    const lang = formData.lang.value;
    updateUserInterface(lang).then(() => {
      loggedUserInfoStore.getState().update({
        ...loggedUserInfo,
        language: lang,
      });
      localize.setupAppLanguage();
      toast.onShow({
        msg: t('update', { keyPrefix: 'toast' }),
        variant: 'success',
      });
    });
  };

  useEffect(() => {
    getLangs();
  }, []);
  return (
    <Form noValidate onSubmit={handleSubmit}>
      <Form.Group controlId="emailSend" className="mb-3">
        <Form.Label>{t('lang.label')}</Form.Label>
        <Form.Select
          value={formData.lang.value}
          isInvalid={formData.lang.isInvalid}
          onChange={(e) => {
            setFormData({
              lang: {
                value: e.target.value,
                isInvalid: false,
                errorMsg: '',
              },
            });
          }}>
          {langs?.map((item) => {
            return (
              <option value={item.value} key={item.label}>
                {item.label}
              </option>
            );
          })}
        </Form.Select>
        <Form.Text as="div">{t('lang.text')}</Form.Text>
        <Form.Control.Feedback type="invalid">
          {formData.lang.errorMsg}
        </Form.Control.Feedback>
      </Form.Group>

      <Button variant="primary" type="submit">
        {t('save', { keyPrefix: 'btns' })}
      </Button>
    </Form>
  );
};

export default React.memo(Index);
