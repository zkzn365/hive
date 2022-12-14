import { FC, memo, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { JSONSchema, SchemaForm, UISchema } from '@/components';
import { FormDataType } from '@/common/interface';
import { brandSetting, getBrandSetting } from '@/services';
import { brandingStore } from '@/stores';
import { useToast } from '@/hooks';
import { handleFormError } from '@/utils';

const uploadType = 'branding';
const Index: FC = () => {
  const { t } = useTranslation('translation', {
    keyPrefix: 'admin.branding',
  });
  const { branding: brandingInfo, update } = brandingStore();
  const Toast = useToast();

  const [formData, setFormData] = useState<FormDataType>({
    logo: {
      value: brandingInfo.logo,
      isInvalid: false,
      errorMsg: '',
    },
    mobile_logo: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    square_icon: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    favicon: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
  });

  const schema: JSONSchema = {
    title: t('page_title'),
    properties: {
      logo: {
        type: 'string',
        title: t('logo.label'),
        description: t('logo.text'),
      },
      mobile_logo: {
        type: 'string',
        title: t('mobile_logo.label'),
        description: t('mobile_logo.text'),
      },
      square_icon: {
        type: 'string',
        title: t('square_icon.label'),
        description: t('square_icon.text'),
      },
      favicon: {
        type: 'string',
        title: t('favicon.label'),
        description: t('favicon.text'),
      },
    },
  };

  const uiSchema: UISchema = {
    logo: {
      'ui:widget': 'upload',
      'ui:options': {
        imageType: uploadType,
      },
    },
    mobile_logo: {
      'ui:widget': 'upload',
      'ui:options': {
        imageType: uploadType,
      },
    },
    square_icon: {
      'ui:widget': 'upload',
      'ui:options': {
        imageType: uploadType,
      },
    },
    favicon: {
      'ui:widget': 'upload',
      'ui:options': {
        acceptType: ',image/x-icon,image/vnd.microsoft.icon',
        imageType: uploadType,
      },
    },
  };

  const handleOnChange = (data) => {
    setFormData(data);
  };

  const onSubmit = () => {
    const params = {
      logo: formData.logo.value,
      mobile_logo: formData.mobile_logo.value,
      square_icon: formData.square_icon.value,
      favicon: formData.favicon.value,
    };
    brandSetting(params)
      .then((res) => {
        console.log(res);
        update(params);
        Toast.onShow({
          msg: t('update', { keyPrefix: 'toast' }),
          variant: 'success',
        });
      })
      .catch((err) => {
        if (err.isError) {
          const data = handleFormError(err, formData);
          setFormData({ ...data });
        }
      });
  };

  const getBrandData = async () => {
    const res = await getBrandSetting();
    if (res) {
      formData.logo.value = res.logo;
      formData.mobile_logo.value = res.mobile_logo;
      formData.square_icon.value = res.square_icon;
      formData.favicon.value = res.favicon;
      setFormData({ ...formData });
    }
  };

  useEffect(() => {
    getBrandData();
  }, []);

  return (
    <div>
      <h3 className="mb-4">{t('page_title')}</h3>
      <SchemaForm
        schema={schema}
        uiSchema={uiSchema}
        formData={formData}
        onSubmit={onSubmit}
        onChange={handleOnChange}
      />
    </div>
  );
};

export default memo(Index);
